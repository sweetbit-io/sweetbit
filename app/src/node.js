import React, { useCallback, useEffect, useState } from 'react';
import { useModal } from 'react-modal-hook';
import Modal from './modal';
import Status from './status';
import Button from './button';
import ManageNode from './manage-node';
import UnlockNode from './unlock-node';

export default function Node({
  api,
  id,
  type,
  name,
  uri,
  enabled,
  onDelete,
  onEnable,
  status,
  onChangeStatus,
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

  useEffect(() => {
    const client = api.subscribeNodeStatus(id);

    client.onmessage = (event) => {
      const payload = JSON.parse(event.data);
      onChangeStatus(id, payload.status);
    };

    return () => {
      client.close();
    };
  }, [api, onChangeStatus, id]);

  const onUnlock = useCallback((password) => {
    const doUnlock = async () => {
      await api.unlockNode(id, password);
    };
    doUnlock();
  }, [id, api]);

  const onSetUp = useCallback(() => {
    const doUnlock = async () => {
      let seed;

      try {
        seed = await api.generateNodeSeed(id);
      } catch (e) {
        console.error(e);
        return;
      }

      let node;

      try {
        node = await api.initNode(id, 'password', seed.mnemonic);
      } catch (e) {
        console.log(e);
        return;
      }

      console.log(node);
    };
    doUnlock();
  }, [id, api]);

  const [connection, setConnection] = useState();

  const [showConnectNodeModal, hideConnectNodeModal] = useModal(({ in: open }) => (
    <Modal open={open} onClose={hideConnectNodeModal}>
      <ManageNode onClose={hideConnectNodeModal} connection={connection} />
    </Modal>
  ), [api, connection]);

  const onConnect = useCallback(() => {
    const doConnect = async () => {
      let connection;

      try {
        connection = await api.generateNodeConnection(id);
      } catch (e) {
        console.error(e);
        return;
      }

      setConnection(connection);
      showConnectNodeModal();
    };
    doConnect();
  }, [id, api, setConnection, showConnectNodeModal]);

  const [showUnlockNodeModal, hideUnlockNodeModal] = useModal(({ in: open }) => (
    <Modal open={open} onClose={hideUnlockNodeModal}>
      <UnlockNode onSetPassword={onUnlock} onCancel={hideUnlockNodeModal} />
    </Modal>
  ), [api]);

  return (
    <div className="node">
      <div className="icon">
        <Status status={enabled && 'green'} />
      </div>
      <div className="label">
        <h1><input type="text" defaultValue={name} /></h1>
        {status === 'stopped' ? (
          <>
            <p>Locally installed node is currently stopped</p>
            <p>
              <button className="link" onClick={deleteNode}>delete node</button>
            </p>
          </>
        ) : status === 'uninitialized' ? (
          <>
            <p>Node needs to be set up with a wallet in order to start</p>
            <p>
              <button className="link" onClick={disableNode}>stop node</button>
            </p>
          </>
        ) : status === 'locked' ? (
          <>
            <p>Node needs to be unlocked with a password in order to start</p>
            <p>
              <button className="link" onClick={disableNode}>stop node</button>
            </p>
          </>
        ) : status === 'started' ? (
          <>
            <p>Node runs locally on your candy dispenser and can be managed remotely</p>
            <p>
              <button className="link" onClick={onConnect}>manage node</button>
            </p>
          </>
        ) : null}
      </div>
      <div className="action">
        {status === 'stopped' ? (
          <Button outline onClick={enableNode}>start</Button>
        ) : status ===  'uninitialized' ? (
          <Button outline onClick={onSetUp}>set up</Button>
        ) : status ===  'locked' ? (
          <Button outline onClick={showUnlockNodeModal}>unlock</Button>
        ) : status ===  'started' ? (
          <Button outline onClick={disableNode}>stop</Button>
        ) : null}
      </div>
      <style jsx>{`
        .node {
          padding: 20px;
          padding-left: 56px;
          position: relative;
          display: flex;
        }

        input {
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

        .label h1 {
          margin: 0;
          font-size: 16px;
          font-weight: 500;
        }

        .label p {
          margin: 5px 0 0;
          color: #555;
        }

        .action {
          flex: 0;
          padding-left: 10px;
          flex: 0 auto;
        }

        button.link {
          background: transparent;
          border: none;
          font-size: inherit;
          color: #5335B8;
          text-decoration: underline;
          font-weight: inherit;
          padding: 0;
        }
      `}</style>
    </div>
  );
}
