import React, { useEffect, useMemo, useCallback, useState } from 'react';
import { useModal } from 'react-modal-hook';
import semver from 'semver';
import css from 'styled-jsx/css';
import Modal from './modal';
import { useDispenserState } from './hooks/state';
import { useNodesState } from './hooks/state';
import Node from './node';
import NoNodes from './no-nodes';
import AddNode from './add-node';
import Status from './status';
import Button from './button';
import Spinner from './spinner';
import Toggle from './toggle';
import publicUrl from './public-url';
import { ReactComponent as DispenserImage } from './dispenser.svg';

const { className, styles } = css.resolve`
  .image {
    width: auto;
    height: 80px;
  }
`;

function App() {
  const [dispenser, setDispenser] = useDispenserState(null);
  const [nodes, setNodes] = useNodesState([]);
  const [release, setRelease] = useState();

  useEffect(() => {
    async function doFetch() {
      const res = await fetch(`${publicUrl}/api/v1/dispenser`);
      const dispenser = await res.json();
      setDispenser(dispenser);
    }
    doFetch();
  }, [setDispenser]);

  useEffect(() => {
    async function doFetch() {
      const res = await fetch('https://api.github.com/repos/sweetbit-io/sweetbit/releases/latest');
      const release = await res.json();
      setRelease(release);
    }
    doFetch();
  }, [setRelease]);

  const availableUpdate = useMemo(() => {
    if (!dispenser || !release) {
      return null;
    }

    const version = semver.clean(release.tag_name);
    if (semver.lte(version, dispenser.version)) {
      return null;
    }

    const asset = release.assets.find(asset => asset.name.indexOf('.mender') >= 0);
    if (!asset) {
      return null;
    }

    const url = asset.browser_download_url;
    const name = release.name;

    return {
      url,
      version,
      name,
    };
  }, [dispenser, release]);

  const setDispenseOnTouch = useCallback((dispenseOnTouch) => {
    async function doSetDispenseOnTouch() {
      const res = await fetch(`${publicUrl}/api/v1/dispenser`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify([{
          op: 'set',
          name: 'dispenseOnTouch',
          value: dispenseOnTouch,
        }]),
      });
      const dispenser = await res.json();
      setDispenser(dispenser);
    }
    doSetDispenseOnTouch();
  }, [setDispenser]);

  const update = useCallback(() => {
    async function onUpdate() {
      const res = await fetch(`${publicUrl}/api/v1/updates`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          url: availableUpdate.url,
        }),
      });
      const update = await res.json();
      setDispenser({
        ...dispenser,
        update,
      });
    }
    onUpdate();
  }, [availableUpdate, dispenser, setDispenser]);

  useEffect(() => {
    async function doFetch() {
      const res = await fetch(`${publicUrl}/api/v1/nodes`);
      const nodes = await res.json();
      setNodes(nodes);
    }
    doFetch();
  }, [setNodes]);

  const deleteNode = useCallback((id) => {
    async function doDelete() {
      await fetch(`${publicUrl}/api/v1/nodes/${id}`, {
        method: 'DELETE',
      });
      setNodes(nodes.filter(node => node.id !== id));
    }
    doDelete();
  }, [nodes, setNodes]);

  const [showAddNodeModal, hideAddNodeModal] = useModal(({ in: open }) => (
    <Modal open={open} onClose={hideAddNodeModal}>
      <AddNode
        onAdd={addNode}
        onCancel={hideAddNodeModal}
      />
    </Modal>
  ), []);

  const addNode = useCallback((data) => {
    async function onAdd() {
      const res = await fetch(`${publicUrl}/api/v1/nodes`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });
      const node = await res.json();
      setNodes([...nodes, node]);
      hideAddNodeModal();
    }
    onAdd();
  }, [nodes, setNodes, hideAddNodeModal]);

  const enableNode = useCallback((id, enabled) => {
    async function onAdd() {
      const res = await fetch(`${publicUrl}/api/v1/nodes/${id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          enabled,
        }),
      });
      const node = await res.json();
      setNodes(nodes.map(n => n.id === node.id ? node : n));
    }
    onAdd();
  }, [nodes, setNodes]);

  return (
    <div>
      <div className="dispenser">
        <div className="header">
          <DispenserImage className={`${className} image`} />
          <h1>{dispenser && dispenser.name}</h1>
        </div>
        <div className="cell">
          <div className="icon">
            <Status status="green" />
          </div>
          <div className="label">
            <h1>Candy dispenser {dispenser && dispenser.state}</h1>
            <p>Your candy dispenser is fully operational</p>
          </div>
        </div>
        {availableUpdate && (
          <div className="cell">
            <div className="icon">
              {/* <Spinner /> */}
            </div>
            <div className="label">
              <h1>Update {availableUpdate.version} available</h1>
              <p>{availableUpdate.name}</p>
            </div>
            <div className="action">
              <Button outline onClick={update}>Update</Button>
            </div>
          </div>
        )}
      </div>
      <div className="pos">
        <div className="cell">
          <div className="icon">
            <Status status="green" />
          </div>
          <div className="label">
            <h1>Point of sales</h1>
            <p>{dispenser && `${dispenser.pos}.onion`}</p>
          </div>
          <div className="action">
            <Button outline>Open</Button>
            <Button outline>⋯</Button>
          </div>
        </div>
        <div className="cell">
          <div className="icon">
          </div>
          <div className="label">
            <h1>Dispense on touch</h1>
            <p>Dispenses candy even without payment by using the touch sensor</p>
          </div>
          <div className="action">
            <Toggle checked={dispenser && dispenser.dispenseOnTouch} onChange={setDispenseOnTouch} />
          </div>
        </div>
      </div>
      <div className="nodes">
        <div className="title">
          <p className="text">Nodes</p>
          <div className="actions">
            {nodes && nodes.length > 0 && (
              <button onClick={showAddNodeModal}>add</button>
            )}
          </div>
        </div>
        <div className="items">
          {nodes && nodes.length === 0 && (
            <div className="node">
              <NoNodes onAdd={showAddNodeModal} />
            </div>
          )}
          {nodes && nodes.map(node => (
            <div className="node" key={node.id}>
              <Node
                id={node.id}
                type={node.type}
                name={node.name}
                enabled={node.enabled}
                onDelete={deleteNode}
                onEnable={enableNode}
              />
            </div>
          ))}
        </div>
      </div>
      <div className="feedback">
        <a href="https://github.com/sweetbit-io/sweetbit/issues/new">How do you like your candy dispenser?</a>
      </div>
      {styles}
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
          text-align: center;
        }

        @media (min-width: 460px) {
          .dispenser {
            padding-top: 50px;
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
          padding-top: 10px;
        }

        .cell {
          border: 1px solid #f1f1f1;
          padding: 20px;
          background: #fff;
          padding-left: 56px;
          position: relative;
          display: flex;
        }

        .cell + .cell {
          border-top: none;
        }

        .cell .icon {
          position: absolute;
          top: 20px;
          left: 20px;
          width: 24px;
          text-align: center;
        }

        .cell .label {
          flex: 1;
        }

        .cell .action {
          flex: 0;
          padding-left: 10px;
          flex: 0 auto;
        }

        .cell h1 {
          margin: 0;
          font-size: 16px;
          font-weight: 500;
        }

        .cell p {
          margin: 5px 0 0;
          color: #555;
        }

        .pos {
          max-width: 460px;
          margin: 20px auto 0;
        }

        @media (min-width: 460px) {
          .cell:first-child {
            border-top-left-radius: 10px;
            border-top-right-radius: 10px;
          }

          .cell:last-child {
            border-bottom-left-radius: 10px;
            border-bottom-right-radius: 10px;
          }
        }

        .pos h1 {
          margin: 0;
          font-size: 16px;
          font-weight: 500;
        }

        .title {
          padding: 20px 20px 5px;
          display: flex;
        }

        .title .text {
          flex: 1;
          margin: 0;
        }

        .title .actions {
          flex: 0;
        }

        .nodes {
          max-width: 460px;
          margin: 0 auto;
        }

        .nodes .items .node {
          border: 1px solid #f1f1f1;
          background: #fff;
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

        .feedback {
          max-width: 460px;
          margin: 0 auto;
          padding: 20px 20px 80px;
          text-align: center;
        }
      `}</style>
      <style jsx global>{`
        * {
          box-sizing: border-box;
        }

        body {
          background-color: #f8f8f8;
          color: #333;
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

        .ReactModal__Body--open {
          overflow: hidden;
        }
      `}</style>
    </div>
  );
}

export default App;
