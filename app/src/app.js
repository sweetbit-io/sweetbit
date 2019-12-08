import React, { useEffect, useCallback } from 'react';
import { useDispenserState } from './hooks/state';
import { useNodesState } from './hooks/state';
import Node from './node';

function App() {
  const [dispenser, setDispenser] = useDispenserState(null);
  const [nodes, setNodes] = useNodesState([]);

  useEffect(() => {
    async function doFetch() {
      const res = await fetch('http://localhost:9000/api/v1/dispenser');
      const dispenser = await res.json();
      setDispenser(dispenser);
    }
    doFetch();
  }, [setDispenser]);

  useEffect(() => {
    async function doFetch() {
      const res = await fetch('http://localhost:9000/api/v1/nodes');
      const nodes = await res.json();
      setNodes(nodes);
    }
    doFetch();
  }, [setNodes]);

  const deleteNode = useCallback((id) => {
    async function doDelete() {
      await fetch(`http://localhost:9000/api/v1/nodes/${id}`, { method: 'DELETE' });
      setNodes(nodes.filter(node => node.id !== id));
    }
    doDelete();
  }, [nodes, setNodes]);

  return (
    <div>
      <div className="dispenser">
        <div className="header">
          <h1>{dispenser && dispenser.name}</h1>
        </div>
        <div className="info">
          Status {dispenser && dispenser.state}
        </div>
        <div className="info">
          Update avaialble
        </div>
      </div>
      <div className="pos">
        <h1>
          <span>PoS {dispenser && dispenser.pos}</span>
        </h1>
      </div>
      <div className="nodes">
        <div className="actions">
          Order, Remove, Add
        </div>
        <div className="items">
          {nodes && nodes.map(node => (
            <div className="node" key={node.id}>
              <Node
                id={node.id}
                name={node.name}
                enabled={node.enabled}
                onDelete={deleteNode}
              />
            </div>
          ))}
        </div>
      </div>
      <div className="networks">
        <div className="actions">
          Order, Remove, Add
        </div>
        <div className="items">
          <div className="network">
            <h1>
              <span>My network</span>
            </h1>
          </div>
          <div className="network">
            <h1>
              <span>My My network</span>
            </h1>
          </div>
        </div>
      </div>
      <div className="feedback">
        How do you like your candy dispenser?
      </div>
      <style jsx>{`
        .dispenser {
          max-width: 460px;
          margin: 0 auto;
        }

        .dispenser div + div {
          border-top: none;
        }

        .dispenser .header {
          background: #804FA0;
          color: white;
          padding: 20px;
        }

        @media (min-width: 460px) {
          .dispenser {
            margin-top: 20px;
          }

          .dispenser div:first-child {
            border-top-left-radius: 10px;
            border-top-right-radius: 10px;
          }

          .dispenser div:last-child {
            border-bottom-left-radius: 10px;
            border-bottom-right-radius: 10px;
          }
        }

        .dispenser .header h1 {
          margin: 0;
          font-size: 24px;
          font-weight: 500;
        }

        .dispenser .info {
          border: 1px solid #f1f1f1;
          padding: 20px;
        }

        .pos {
          max-width: 460px;
          margin: 20px auto 0;
          border: 1px solid #f1f1f1;
          padding: 20px;
        }

        @media (min-width: 460px) {
          .pos {
            border-radius: 10px;
          }
        }

        .pos h1 {
          margin: 0;
          font-size: 16px;
          font-weight: 500;
        }

        .nodes {
          max-width: 460px;
          margin: 0 auto;
        }

        .nodes .actions {
          padding: 20px 20px 5px;
        }

        .nodes .items .node {
          border: 1px solid #f1f1f1;
          padding: 20px;
        }

        .nodes .items .node + .node {
          border-top: none;
        }

        @media (min-width: 460px) {
          .nodes .items .node:first-child {
            border-top-left-radius: 10px;
            border-top-right-radius: 10px;
          }

          .nodes .items .node:last-child {
            border-bottom-left-radius: 10px;
            border-bottom-right-radius: 10px;
          }
        }

        .networks {
          max-width: 460px;
          margin: 0 auto;
        }

        .networks .actions {
          padding: 20px 20px 5px;
        }

        .networks .items .network {
          border: 1px solid #f1f1f1;
          padding: 20px;
        }

        .networks .items .network + .network {
          border-top: none;
        }

        @media (min-width: 460px) {
          .networks .items .network:first-child {
            border-top-left-radius: 10px;
            border-top-right-radius: 10px;
          }

          .networks .items .network:last-child {
            border-bottom-left-radius: 10px;
            border-bottom-right-radius: 10px;
          }
        }

        .networks .network h1 {
          margin: 0;
          font-size: 16px;
          font-weight: 500;
        }

        .feedback {
          max-width: 460px;
          margin: 0 auto;
          padding: 20px;
          text-align: center;
        }
      `}</style>
      <style jsx global>{`
        * {
          box-sizing: border-box;
        }

        body {
          background-color: #fff;
          color: #222;
          font-size: 16px;
          margin: 0;
          padding: 0;
          font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
          -webkit-font-smoothing: antialiased;
          -moz-osx-font-smoothing: grayscale;
        }

        code {
          font-family: source-code-pro, Menlo, Monaco, Consolas, "Courier New", monospace;
        }
      `}</style>
    </div>
  );
}

export default App;
