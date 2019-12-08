import React, { useState, useEffect } from 'react';
import { useModal } from 'react-modal-hook';
import Node from '../node';
import Modal from '../modal';
import { useNodesState } from '../hooks/state';

function Nodes() {
  const [nodes, setNodes] = useNodesState([]);
  const [count, setCount] = useState(0);

  useEffect(() => {
    async function doFetch() {
      const res = await fetch('http://localhost:9000/api/v1/nodes');
      const nodes = await res.json();
      setNodes(nodes);
    }

    doFetch();
  }, [setNodes]);

  const [showModal] = useModal(({ in: open, onExited }) => (
    <Modal open={open} onClose={onExited}>
      <span>The count is {count}</span>
      <button onClick={() => setCount(count + 1)}>Increment</button>
    </Modal>
  ), [count]);

  function rename() {
  }

  function unlock() {
    showModal();
  }

  function deleteNode() {
  }

  return (
    <div>
      {nodes.map(node => (
        <Node
          id={node.id}
          name={node.name}
          onRename={rename}
          status="locked"
          onUnlock={unlock}
          onDelete={deleteNode}
        />
      ))}
      <div className="">
        add new...
      </div>
      <style jsx>{`
        article {
          display: block;
          background: white;
          padding: 10px;
        }
      `}</style>
    </div>
  );
}

export default Nodes;
