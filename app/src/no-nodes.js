import React, { useCallback } from 'react';
import css from 'styled-jsx/css';
import { ReactComponent as NoNodesImage } from './no-nodes.svg';
import Button from './button';

const { className, styles } = css.resolve`
  .no-nodes {
    width: auto;
    height: 120px;
  }
`;

export default function NoNodes({
  onAdd,
}) {
  const addNode = useCallback((e) => {
    onAdd();
  }, [onAdd]);

  return (
    <div>
      <NoNodesImage className={`${className} no-nodes`} />
      <p>Add a Lightning node through which you can accept payments.</p>
      <Button onClick={addNode}>add node</Button>
      {styles}
      <style jsx>{`
        div {
          text-align: center;
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
