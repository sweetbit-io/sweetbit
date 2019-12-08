import React, { useCallback } from 'react';

function Node({
  id,
  name,
  onRename,
  onDelete,
  status,
  onUnlock,
}) {
  const deleteNode = useCallback(() => {
    async function doDelete() {
      await fetch(`http://localhost:9000/api/v1/nodes/${id}`, { method: 'DELETE' });
      onDelete();
    }
    doDelete();
  }, [id, onDelete]);

  const renameNode = useCallback(() => {
    async function doRename() {
      await fetch(`http://localhost:9000/api/v1/nodes/${id}`, {
        method: 'PATCH',
        headers: {
          'content-type': 'application/json',
        },
        body: JSON.stringify(),
      });
    }
    doRename();
  }, [id]);

  return (
    <article>
      <div>remote node</div>
      <h1>{name}</h1>
      <p>
        <button onClick={onUnlock}>{status}</button>
      </p>
      <button onClick={deleteNode}>delete</button>
      <style jsx>{`
        article {
          margin: 10px;
          border-radius: 5px;
          box-shadow: 0px 1px 5px rgba(0,0,0,0.2);
          padding-left: 40px;
          display: block;
          background: white;
          padding: 10px;
        }
      `}</style>
    </article>
  );
}

export default Node;
