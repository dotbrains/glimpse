document.addEventListener('DOMContentLoaded', async () => {
  const resp = await fetch('/api/diff');
  const data = await resp.json();

  document.getElementById('repo-name').textContent = data.repo || '';
  const refInfo = document.getElementById('ref-info');
  if (data.base && data.compare) {
    refInfo.textContent = data.base + '..' + data.compare;
  } else if (data.base) {
    refInfo.textContent = data.base;
  } else {
    refInfo.textContent = 'working tree';
  }
  document.getElementById('summary').textContent = data.summary || '';
  document.title = 'glimpse — ' + (data.repo || 'diff');

  const nav = document.getElementById('file-list');
  const container = document.getElementById('diff-container');

  if (!data.files || data.files.length === 0) {
    container.innerHTML = '<div class="empty">No changes found.</div>';
    return;
  }

  // Build file nav.
  data.files.forEach((file, idx) => {
    const a = document.createElement('a');
    a.href = '#file-' + idx;
    const name = file.newName || file.oldName;
    const stats = fileStats(file);
    a.innerHTML = esc(name) + ' <span class="file-stat"><span class="add">+' + stats.add + '</span> <span class="del">-' + stats.del + '</span></span>';
    nav.appendChild(a);
  });

  // Render each file.
  data.files.forEach((file, idx) => {
    const card = document.createElement('div');
    card.className = 'file-card';
    card.id = 'file-' + idx;

    const name = file.newName || file.oldName;
    const stats = fileStats(file);
    let statusBadge = '';
    if (file.status && file.status !== 'modified') {
      statusBadge = '<span class="status-badge ' + file.status + '">' + file.status + '</span>';
    }

    card.innerHTML =
      '<div class="file-header">' +
        '<span class="path">' + statusBadge + esc(name) + '</span>' +
        '<span class="stats"><span class="add">+' + stats.add + '</span><span class="del">-' + stats.del + '</span></span>' +
      '</div>';

    const table = document.createElement('table');
    table.className = 'diff-table';

    (file.hunks || []).forEach(hunk => {
      // Hunk header row.
      const hr = document.createElement('tr');
      hr.className = 'hunk-header';
      hr.innerHTML = '<td class="line-num"></td><td class="line-num"></td><td class="line-content" colspan="1">' + esc(hunk.header) + '</td>';
      table.appendChild(hr);

      (hunk.lines || []).forEach(line => {
        const tr = document.createElement('tr');
        tr.className = line.type;

        const oldNum = line.type === 'added' ? '' : (line.oldNum || '');
        const newNum = line.type === 'removed' ? '' : (line.newNum || '');
        const prefix = line.type === 'added' ? '+' : line.type === 'removed' ? '-' : ' ';

        tr.innerHTML =
          '<td class="line-num">' + oldNum + '</td>' +
          '<td class="line-num">' + newNum + '</td>' +
          '<td class="line-content">' + prefix + esc(line.content) + '</td>';
        table.appendChild(tr);
      });
    });

    card.appendChild(table);
    container.appendChild(card);
  });
});

function fileStats(file) {
  let add = 0, del = 0;
  (file.hunks || []).forEach(h => {
    (h.lines || []).forEach(l => {
      if (l.type === 'added') add++;
      if (l.type === 'removed') del++;
    });
  });
  return { add, del };
}

function esc(s) {
  const d = document.createElement('div');
  d.textContent = s || '';
  return d.innerHTML;
}
