import React from 'react';
import Markdown from 'react-markdown';
import Button from './button';

export default function Update({
  onCancel,
  onUpdate,
  body,
}) {
  return (
    <div className="update">
      <h1>Update</h1>
      <Markdown source={body} />
      <div className="actions">
        <div className="action">
          <Button type="button" onClick={onUpdate}>update</Button>
        </div>
        <div className="action">
          <Button type="button" onClick={onCancel} outline>cancel</Button>
        </div>
      </div>
      <style jsx>{`
        .update {
          padding: 20px;
          background: #fff;
          border-radius: 10px;
        }

        .actions {
          padding-top: 40px;
          display: flex;
        }

        .action + .action {
          padding-left: 10px;
        }

        h1 {
          margin: 0;
        }

        .update :global(img) {
          max-width: 100%;
          width: 100%;
          height: auto;
        }
      `}</style>
    </div>
  );
}
