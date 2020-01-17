import React, { useMemo } from 'react';
import css from 'styled-jsx/css';
import { useModal } from 'react-modal-hook';
import QRCode from 'qrcode.react';
import Modal from './modal';
import Button from './button';
import { ReactComponent as ManageNodeImage } from './manage-node.svg';

const { className: imageClassName, styles: imageStyles } = css.resolve`
  .image {
    width: auto;
    height: 120px;
  }
`;

const { className, styles } = css.resolve`
  svg {
    width: 100%;
    height: auto;
  }
`;

export default function ManageNode({
  connection,
  onClose,
}) {
  const lndconnect = useMemo(() => {
    if (!connection) {
      return null;
    }

    return `lndconnect://${connection.uri}?cert=${encodeURIComponent(connection.cert)}&macaroon=${encodeURIComponent(connection.macaroon)}`;
  }, [connection]);

  const [showQRModal, hideQRModal] = useModal(({ in: open }) => (
    <Modal open={open} onClose={hideQRModal}>
      <div className="lndconnect">
        <p>Scan this QR code when connecting to a remote node in Zap.</p>
        <QRCode className={className} renderAs="svg" size={512} value={JSON.stringify({
          ip: connection.uri,
          c: connection.cert,
          m: connection.macaroon,
        })} />
        <p>
          <Button outline onClick={hideQRModal}>close</Button>
        </p>
      </div>
      {styles}
      <style jsx>{`
        .lndconnect {
          padding: 20px;
          text-align: center;
        }
      `}</style>
    </Modal>
  ), [connection]);

  return (
    <div className="manage-node">
      <p className="center">
        <ManageNodeImage className={`${imageClassName} image`} />
      </p>
      <h1 className="center">Manage your node</h1>
      <p className="center">Manage your node's funds and channels through Zap.</p>
      <p className="center">
        <Button href={lndconnect}>
          open in Zap
        </Button>
      </p>
      <p className="center">
        or <button className="link" onClick={showQRModal}>display QR code</button>
      </p>
      <p className="center">
        Zap is available on Android, iOS and on desktop.
        It <a href="https://zap.jackmallers.com" target="_blank" rel="noopener noreferrer">can be downloaded from its website</a>.
      </p>
      <div className="center actions">
        <Button outline onClick={onClose}>close</Button>
      </div>
      {imageStyles}
      <style jsx>{`
        .manage-node {
          padding: 20px;
        }

        .center {
          text-align: center;
        }

        h1 {
          margin: 0;
        }

        .link {
          background: transparent;
          border: none;
          font-size: inherit;
          color: #5335B8;
          text-decoration: underline;
          font-weight: inherit;
          padding: 0;
        }

        .actions {
          padding-top: 40px;
        }
      `}</style>
    </div>
  );
}
