let allComments = [];

document.addEventListener('DOMContentLoaded', async () => {
  const [diffResp, commentsResp] = await Promise.all([
    fetch('/api/diff'),
    fetch('/api/comments')
  ]);
  const data = await diffResp.json();
  allComments = await commentsResp.json() || [];

  document.getElementById('repo-name').textContent = data.repo || '';
  const refInfo = document.getElementById('ref-info');
  if (data.base && data.compare) refInfo.textContent = data.base + '..' + data.compare;
  else if (data.base) refInfo.textContent = data.base;
  else refInfo.textContent = 'working tree';
  document.getElementById('summary').textContent = data.summary || '';
  document.title = 'glimpse — ' + (data.repo || 'diff');

  const nav = document.getElementById('file-list');
  const container = document.getElementById('diff-container');

  if (!data.files || data.files.length === 0) {
    container.innerHTML = '<div class="empty">No changes found.</div>';
    return;
  }

  data.files.forEach((file, idx) => {
    const a = document.createElement('a');
    a.href = '#file-' + idx;
    const name = file.newName || file.oldName;
    const stats = fileStats(file);
    a.innerHTML = esc(name) + ' <span class="file-stat"><span class="add">+' + stats.add + '</span> <span class="del">-' + stats.del + '</span></span>';
    nav.appendChild(a);
  });

  data.files.forEach((file, idx) => {
    const card = document.createElement('div');
    card.className = 'file-card';
    card.id = 'file-' + idx;
    const name = file.newName || file.oldName;
    const stats = fileStats(file);
    let statusBadge = '';
    if (file.status && file.status !== 'modified') statusBadge = '<span class="status-badge ' + file.status + '">' + file.status + '</span>';
    card.innerHTML = '<div class="file-header"><span class="path">' + statusBadge + esc(name) + '</span><span class="stats"><span class="add">+' + stats.add + '</span><span class="del">-' + stats.del + '</span></span></div>';

    const table = document.createElement('table');
    table.className = 'diff-table';

    (file.hunks || []).forEach(hunk => {
      const hr = document.createElement('tr');
      hr.className = 'hunk-header';
      hr.innerHTML = '<td class="line-num"></td><td class="line-num"></td><td class="line-content">' + esc(hunk.header) + '</td>';
      table.appendChild(hr);

      (hunk.lines || []).forEach(line => {
        const tr = document.createElement('tr');
        tr.className = line.type;
        const oldNum = line.type === 'added' ? '' : (line.oldNum || '');
        const newNum = line.type === 'removed' ? '' : (line.newNum || '');
        const prefix = line.type === 'added' ? '+' : line.type === 'removed' ? '-' : ' ';
        tr.innerHTML = '<td class="line-num">' + oldNum + '</td><td class="line-num comment-gutter" data-file="' + esc(name) + '" data-line="' + (newNum || oldNum) + '" data-side="' + (line.type === 'removed' ? 'old' : 'new') + '">' + newNum + '</td><td class="line-content">' + prefix + esc(line.content) + '</td>';
        table.appendChild(tr);

        // Render inline comments for this line.
        const lineComments = allComments.filter(c => c.file === name && c.line === (newNum || oldNum));
        lineComments.forEach(c => {
          const cr = document.createElement('tr');
          cr.className = 'comment-row' + (c.resolved ? ' resolved' : '');
          cr.innerHTML = '<td class="line-num"></td><td class="line-num"></td><td class="comment-cell"><span class="severity-tag ' + c.severity + '">' + c.severity + '</span> <span class="comment-body">' + esc(c.body) + '</span> <span class="comment-author">' + c.author + '</span>' + (!c.resolved ? ' <button class="resolve-btn" data-id="' + c.id + '">resolve</button>' : ' <span class="resolved-tag">resolved</span>') + '</td>';
          table.appendChild(cr);
        });
      });
    });

    card.appendChild(table);
    container.appendChild(card);
  });

  // Click gutter to add comment.
  document.addEventListener('click', async (e) => {
    const gutter = e.target.closest('.comment-gutter');
    if (!gutter) return;
    const file = gutter.dataset.file;
    const line = parseInt(gutter.dataset.line);
    const side = gutter.dataset.side;
    if (!file || !line) return;
    const body = prompt('Comment:');
    if (!body) return;
    const sev = prompt('Severity (must-fix, suggestion, nit, question):', 'suggestion') || 'suggestion';
    await fetch('/api/comments', { method: 'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify({file, line, side, body, severity: sev}) });
    location.reload();
  });

  // Resolve button.
  document.addEventListener('click', async (e) => {
    if (!e.target.classList.contains('resolve-btn')) return;
    const id = e.target.dataset.id;
    await fetch('/api/comments/' + id, { method: 'PATCH' });
    location.reload();
  });
});

function fileStats(file) {
  let add = 0, del = 0;
  (file.hunks || []).forEach(h => (h.lines || []).forEach(l => { if (l.type === 'added') add++; if (l.type === 'removed') del++; }));
  return { add, del };
}

function esc(s) { const d = document.createElement('div'); d.textContent = s || ''; return d.innerHTML; }
