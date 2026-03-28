document.addEventListener('DOMContentLoaded', async () => {
  const [treeResp, commentsResp] = await Promise.all([
    fetch('/api/tree'),
    fetch('/api/comments')
  ]);
  const treeData = await treeResp.json();
  const allComments = await commentsResp.json() || [];

  document.getElementById('repo-name').textContent = treeData.repo || '';
  document.getElementById('file-count').textContent = (treeData.files || []).length + ' files';
  document.title = 'glimpse — ' + (treeData.repo || 'tree');

  const sidebar = document.getElementById('tree-sidebar');
  const viewer = document.getElementById('file-viewer');

  // Build tree structure from flat file list.
  const tree = {};
  (treeData.files || []).forEach(path => {
    const parts = path.split('/');
    let node = tree;
    parts.forEach((part, i) => {
      if (!node[part]) node[part] = i === parts.length - 1 ? null : {};
      if (node[part] !== null) node = node[part];
    });
  });

  function renderTree(node, prefix, container) {
    const entries = Object.keys(node).sort((a, b) => {
      const aDir = node[a] !== null;
      const bDir = node[b] !== null;
      if (aDir !== bDir) return aDir ? -1 : 1;
      return a.localeCompare(b);
    });
    entries.forEach(name => {
      const fullPath = prefix ? prefix + '/' + name : name;
      const isDir = node[name] !== null;
      const el = document.createElement('div');
      el.className = isDir ? 'tree-dir' : 'tree-file';
      el.textContent = (isDir ? '📁 ' : '  ') + name;
      el.dataset.path = fullPath;
      container.appendChild(el);
      if (isDir) {
        const children = document.createElement('div');
        children.className = 'tree-children';
        children.style.display = 'none';
        renderTree(node[name], fullPath, children);
        container.appendChild(children);
        el.addEventListener('click', (e) => {
          e.stopPropagation();
          children.style.display = children.style.display === 'none' ? 'block' : 'none';
          el.textContent = (children.style.display === 'none' ? '📁 ' : '📂 ') + name;
        });
      } else {
        el.addEventListener('click', () => loadFile(fullPath));
      }
    });
  }
  renderTree(tree, '', sidebar);

  async function loadFile(path) {
    const resp = await fetch('/api/file?path=' + encodeURIComponent(path));
    if (!resp.ok) { viewer.innerHTML = '<div class="empty">Could not load file.</div>'; return; }
    const data = await resp.json();
    const lines = (data.content || '').split('\n');
    const fileComments = allComments.filter(c => c.file === path);

    let html = '<div class="file-card"><div class="file-header"><span class="path">' + esc(path) + '</span></div>';
    html += '<table class="diff-table">';
    lines.forEach((line, i) => {
      const num = i + 1;
      html += '<tr class="context"><td class="line-num comment-gutter" data-file="' + esc(path) + '" data-line="' + num + '" data-side="new">' + num + '</td><td class="line-content"> ' + esc(line) + '</td></tr>';
      const lc = fileComments.filter(c => c.line === num);
      lc.forEach(c => {
        html += '<tr class="comment-row' + (c.resolved ? ' resolved' : '') + '"><td class="line-num"></td><td class="comment-cell"><span class="severity-tag ' + c.severity + '">' + c.severity + '</span> ' + esc(c.body) + ' <span class="comment-author">' + c.author + '</span>' + (!c.resolved ? ' <button class="resolve-btn" data-id="' + c.id + '">resolve</button>' : '') + '</td></tr>';
      });
    });
    html += '</table></div>';
    viewer.innerHTML = html;
  }

  // Click gutter to comment.
  document.addEventListener('click', async (e) => {
    const gutter = e.target.closest('.comment-gutter');
    if (!gutter) return;
    const file = gutter.dataset.file, line = parseInt(gutter.dataset.line);
    if (!file || !line) return;
    const body = prompt('Comment:');
    if (!body) return;
    const sev = prompt('Severity (must-fix, suggestion, nit, question):', 'suggestion') || 'suggestion';
    await fetch('/api/comments', { method: 'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify({file, line, side: 'new', body, severity: sev}) });
    location.reload();
  });

  document.addEventListener('click', async (e) => {
    if (!e.target.classList.contains('resolve-btn')) return;
    await fetch('/api/comments/' + e.target.dataset.id, { method: 'PATCH' });
    location.reload();
  });
});

function esc(s) { const d = document.createElement('div'); d.textContent = s || ''; return d.innerHTML; }
