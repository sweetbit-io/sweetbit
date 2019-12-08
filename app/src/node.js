import React, { useCallback } from 'react';

export default function Node({
  id,
  type,
  name,
  enabled,
  onDelete,
  onEnable,
}) {
  const deleteNode = useCallback(() => {
    onDelete(id);
  }, [id, onDelete]);

  const enableNode = useCallback(() => {
    onEnable(id);
  }, [id, onEnable]);

  return (
    <div>
      <h1>
        <span>{name}</span>
        <button onClick={enableNode}>enable</button>
        <button onClick={deleteNode}>delete</button>
      </h1>
      <div className="status">
        {type} {enabled && 'enabled'}
      </div>
      <style jsx>{`
        div {
          padding: 20px;
        }

        h1 {
          margin: 0;
          font-size: 16px;
          font-weight: 500;
        }
      `}</style>
    </div>
  );
}
