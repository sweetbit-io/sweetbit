import React, { useEffect, useMemo, useCallback, useState, useReducer } from 'react';
import { useModal } from 'react-modal-hook';
import semver from 'semver';
import css from 'styled-jsx/css';
import Modal from './modal';
import { useDispenserState } from './hooks/state';
import Node from './node';
import NoNodes from './no-nodes';
import AddNode from './add-node';
import Update from './update';
import Status from './status';
import Button from './button';
import Progress from './progress';
import Toggle from './toggle';
import { publicUrl, publicWsUrl } from './config';
import { useApi } from './api';
import { ReactComponent as DispenserImage } from './dispenser.svg';

const { className, styles } = css.resolve`
  .image {
    width: auto;
    height: 80px;
  }
`;

function nodesReducer(state, action) {
  switch (action.type) {
    case 'set':
      return action.nodes;
    case 'add':
      return [...state, action.node];
    case 'update':
      return state.map(node => node.id === action.node.id ? action.node : node);
    case 'remove':
      return state.filter(node => node.id !== action.id);
    case 'status':
      return state.map(node => node.id === action.id ? { ...node, status: action.status } : node);
    default:
      throw new Error();
  }
}

function App() {
  const api = useApi({ publicUrl, publicWsUrl });
  const [dispenser, setDispenser] = useDispenserState(null);
  // const [nodes, setNodes] = useNodesState([]);
  const [nodes, dispatchNodesAction] = useReducer(nodesReducer, []);
  const [release, setRelease] = useState();
  const [currentUpdate, setCurrentUpdate] = useState();

  useEffect(() => {
    async function doFetch() {
      const res = await fetch(`${publicUrl}/api/v1/dispenser`);
      const dispenser = await res.json();
      setDispenser(dispenser);
    }
    doFetch();
  }, [setDispenser]);

  const dismissUpdate = useCallback(() =>  {
    setDispenser({
      ...dispenser,
      update: null,
    });
  }, [dispenser, setDispenser]);

  const currentUpdateId = useMemo(() => dispenser && dispenser.update && dispenser.update.id, [dispenser]);

  useEffect(() => {
    async function doFetch() {
      if (!currentUpdateId) {
        setCurrentUpdate(null);
        return;
      }
      const res = await fetch(`${publicUrl}/api/v1/updates/${currentUpdateId}`);
      const update = await res.json();
      setCurrentUpdate(update);
    }
    doFetch();
  }, [currentUpdateId, setCurrentUpdate]);

  useEffect(() => {
    if (!currentUpdateId) {
      return;
    }

    const socket = new WebSocket(`${publicWsUrl}/api/v1/updates/${currentUpdateId}/events`);
    socket.onmessage = (event) => {
      const update = JSON.parse(event.data);
      setCurrentUpdate(update);
    };

    return () => {
      socket.close();
    };
  }, [currentUpdateId, setCurrentUpdate]);

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
    const body = release.body;
    const description = release.body.split('.\r\n')[0];

    return {
      url,
      version,
      name,
      description,
      body,
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

  const reboot = useCallback(() => {
    async function doReboot() {
      const res = await fetch(`${publicUrl}/api/v1/dispenser`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify([{
          op: 'reboot',
        }]),
      });
      const dispenser = await res.json();
      setDispenser(dispenser);
    }
    doReboot();
  }, [setDispenser]);

  const [showUpdateModal, hideUpdateModal] = useModal(({ in: open }) => {
    return (
      <Modal open={open} onClose={hideUpdateModal}>
        <Update
          name={release.name}
          body={release.body}
          onUpdate={update}
          onCancel={hideUpdateModal}
        />
      </Modal>
    );
  }, [release]);

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
      hideUpdateModal();
      setDispenser({
        ...dispenser,
        update: { id: update.id },
      });

    }
    onUpdate();
  }, [availableUpdate, dispenser, setDispenser, hideUpdateModal]);

  const cancelUpdate = useCallback(() => {
    async function onCancel() {
      if (!currentUpdate) {
        return;
      }

      const res = await fetch(`${publicUrl}/api/v1/updates/${currentUpdate.id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          state: 'cancelled',
        }),
      });
      const update = await res.json();
      setCurrentUpdate(update);
    }
    onCancel();
  }, [currentUpdate, setCurrentUpdate]);

  const completeUpdate = useCallback(() => {
    async function onComplete() {
      if (!currentUpdate) {
        return;
      }

      const res = await fetch(`${publicUrl}/api/v1/updates/${currentUpdate.id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          state: 'completed',
        }),
      });
      const update = await res.json();
      setCurrentUpdate(update);
    }
    onComplete();
  }, [currentUpdate, setCurrentUpdate]);

  useEffect(() => {
    async function doFetch() {
      const res = await fetch(`${publicUrl}/api/v1/nodes`);
      const nodes = await res.json();
      dispatchNodesAction({
        type: 'set',
        nodes,
      });
    }
    doFetch();
  }, [dispatchNodesAction]);

  const deleteNode = useCallback((id) => {
    async function doDelete() {
      await fetch(`${publicUrl}/api/v1/nodes/${id}`, {
        method: 'DELETE',
      });
      dispatchNodesAction({
        type: 'remove',
        id,
      });
    }
    doDelete();
  }, [dispatchNodesAction]);

  const [showAddNodeModal, hideAddNodeModal] = useModal(({ in: open }) => {
    return (
      <Modal open={open} onClose={hideAddNodeModal}>
        <AddNode
          api={api}
          dispatchNodesAction={dispatchNodesAction}
          onCancel={hideAddNodeModal}
        />
      </Modal>
    );
  }, [api]);

  const enableNode = useCallback((id, enabled) => {
    async function onAdd() {
      const node = await api.enableNode(id, enabled);
      dispatchNodesAction({
        type: 'update',
        node,
      });
    }
    onAdd();
  }, [api, dispatchNodesAction]);

  const changeNodeStatus = useCallback((id, status) => {
    dispatchNodesAction({
      type: 'status',
      id,
      status,
    });
  }, [dispatchNodesAction]);

  return (
    <div>
      <div className="dispenser">
        <div className="header">
          <DispenserImage className={`${className} image`} />
          <h1><input type="text" defaultValue={dispenser && dispenser.name} /></h1>
        </div>
        <div className="cell">
          <div className="icon">
            {!dispenser ? (
              <Status />
            ) : dispenser.state === 'running' ? (
              <Status status="green" />
            ) : dispenser.state ===  'stopping' ? (
              <Status status="red" />
            ) : null}
          </div>
          <div className="label">
            <h1>Candy dispenser {dispenser && dispenser.state}</h1>
            {!dispenser ? (
              <p>Loading...</p>
            ) : dispenser.state === 'running' ? (
              <p>Your candy dispenser is fully operational</p>
            ) : dispenser.state === 'stopping' ? (
              <p>Your candy dispenser is shutting down...</p>
            ) : null}
          </div>
        </div>
        {(availableUpdate || currentUpdate) && (
          <div className="cell">
            <div className="icon">
              {!currentUpdate ? (
                <span role="img" aria-label="refresh">🔄</span>
              ) : currentUpdate.state === 'started' ? (
                <Progress value={currentUpdate.progress} />
              ) : currentUpdate.state === 'cancelled' ? (
                <span role="img" aria-label="cancelled">🚫</span>
              ) : currentUpdate.state === 'failed' ? (
                <span role="img" aria-label="failed">❌</span>
              ) : currentUpdate.state === 'installed' ? (
                <span role="img" aria-label="check">✅</span>
              ) : currentUpdate.state === 'rejected' ? (
                <span role="img" aria-label="rejected">🚫</span>
              ) : currentUpdate.state === 'completed' ? (
                <span role="img" aria-label="check">✅</span>
              ) : null}
            </div>
            <div className="label">
              {!currentUpdate ? (
                <h1>Update {availableUpdate && availableUpdate.version} available</h1>
              ) : currentUpdate.state === 'started' ? (
                <h1>Updating to {availableUpdate && availableUpdate.version}...</h1>
              ) : currentUpdate.state === 'cancelled' ? (
                <h1>Cancelled update</h1>
              ) : currentUpdate.state === 'failed' ? (
                <h1>Update failed</h1>
              ) : currentUpdate.state === 'installed' && currentUpdate.reboot ? (
                <h1>Reboot to complete installation</h1>
              ) : currentUpdate.state === 'installed' && currentUpdate.commit ? (
                <h1>Confirm the installation to complete the update</h1>
              ) : currentUpdate.state === 'rejected' ? (
                <h1>Rejected update</h1>
              ) : currentUpdate.state === 'completed' ? (
                <h1>Successfully completed update</h1>
              ) : null}
              <p>{availableUpdate && availableUpdate.description}</p>
            </div>
            <div className="action">
              {!currentUpdate ? (
                <Button outline onClick={showUpdateModal}>update</Button>
              ) : currentUpdate.state === 'started' ? (
                <Button outline onClick={cancelUpdate}>cancel</Button>
              ) : currentUpdate.state === 'cancelled' ? (
                <Button outline onClick={dismissUpdate}>dismiss</Button>
              ) : currentUpdate.state === 'failed' ? (
                <Button outline onClick={dismissUpdate}>dismiss</Button>
              ) : currentUpdate.state === 'installed' && currentUpdate.reboot ? (
                <Button outline onClick={reboot}>reboot</Button>
              ) : currentUpdate.state === 'installed' && currentUpdate.commit ? (
                <Button outline onClick={completeUpdate}>confirm</Button>
              ) : currentUpdate.state === 'rejected' ? (
                <Button outline onClick={dismissUpdate}>dismiss</Button>
              ) : currentUpdate.state === 'completed' ? (
                <Button outline onClick={dismissUpdate}>dismiss</Button>
              ) : null}
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
            <p>Generates invoices that dispense candy when paid</p>
            {dispenser ? (
              <p>
                <a
                  href={`http://${dispenser.pos}.onion`}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {`http://${dispenser.pos}.onion`}
                </a>
              </p>
            ) : null}
          </div>
          <div className="action">
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
              <button className="action" onClick={showAddNodeModal}>add node</button>
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
                api={api}
                id={node.id}
                type={node.type}
                name={node.name}
                uri={node.uri}
                enabled={node.enabled}
                onDelete={deleteNode}
                onEnable={enableNode}
                status={node.status}
                onChangeStatus={changeNodeStatus}
              />
            </div>
          ))}
        </div>
      </div>
      <div className="feedback">
        <p className="meta">{dispenser && `v${dispenser.version}`}</p>
        <p><a href="https://github.com/sweetbit-io/sweetbit/issues/new">How do you like your candy dispenser?</a></p>
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

        .dispenser input {
          appearance: none;
          border: none;
          font-size: inherit;
          font-weight: inherit;
          font-style: inherit;
          font-kerning: inherit;
          color: inherit;
          width: 100%;
          padding: 0;
          outline: none;
          background: transparent;
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
          overflow: auto;
          flex: 1 1 auto;
        }

        .cell .action {
          flex: 1 0 auto;
          padding-left: 10px;
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
          flex: 0 0 auto;
        }

        .title .actions .action {
          background: transparent;
          border: none;
          font-size: inherit;
          color: #5335B8;
          text-decoration: underline;
          font-weight: inherit;
          padding: 0;
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

        .feedback p {
          margin: 5px 0 0;
        }

        .feedback .meta {
          color: #666;
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

        code {
          overflow-wrap: break-word;
        }

        .ReactModal__Body--open {
          overflow: hidden;
        }
      `}</style>
    </div>
  );
}

export default App;
