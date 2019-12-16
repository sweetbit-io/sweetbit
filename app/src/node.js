import React, { useCallback } from 'react';
import Status from './status';
import Button from './button';

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
    onEnable(id, true);
  }, [id, onEnable]);

  const disableNode = useCallback(() => {
    onEnable(id, false);
  }, [id, onEnable]);

  return (
    <div className="node">
      <div className="icon">
        <Status status={enabled && 'green'} />
      </div>
      <div className="label">
        <h1>{name}</h1>
        <p>{type}</p>
      </div>
      <div className="action">
        {enabled ? (
          <Button outline onClick={disableNode}>disable</Button>
        ) : (
          <Button outline onClick={enableNode}>enable</Button>
        )}
        <Button outline onClick={deleteNode}>⋯</Button>
      </div>
      <style jsx>{`
        .node {
          padding: 20px;
          padding-left: 56px;
          position: relative;
          display: flex;
        }

        .icon {
          position: absolute;
          top: 20px;
          left: 20px;
          width: 24px;
          text-align: center;
        }

        h1 {
          margin: 0;
          font-size: 16px;
          font-weight: 500;
        }

        .label {
          flex: 1;
        }

        .action {
          flex: 0;
          padding-left: 10px;
          flex: 0 auto;
        }
      `}</style>
    </div>
  );
}
