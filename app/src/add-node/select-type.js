import React, { useCallback } from 'react';
import css from 'styled-jsx/css';
import Button from '../button';
import { ReactComponent as NoNodesImage } from '../no-nodes.svg';

const { className, styles } = css.resolve`
  .image {
    width: auto;
    height: 120px;
  }
`;

export default function SelectType({
  onSelect,
  onCancel,
}) {
  const selectLocal = useCallback(() => {
    onSelect('local');
  }, [onSelect]);

  const selectRemoteLND = useCallback(() => {
    onSelect('remote-lnd');
  }, [onSelect]);

  return (
    <div className="select-type">
      <p className="center">
        <NoNodesImage className={`${className} image`} />
      </p>
      <h1 className="center">Select node type</h1>
      <p className="center">Please choose which type of node you'd like to set up.</p>
      <ul className="items">
        <li className="item">
          <button onClick={selectLocal} className="option">
            <strong className="title">Local Node</strong>
            <span className="separator"> – </span>
            <span className="description">Run a Lightning node on the candy dispenser.</span>
            <span className="separator"> (</span>
            <em className="label">recommended</em>
            <span className="separator">)</span>
          </button>
        </li>
        <li className="item">
          <button onClick={selectRemoteLND} className="option">
            <strong className="title">Remote LND Node</strong>
            <span className="separator"> – </span>
            <span className="description">Connect to your own external LND node.</span>
          </button>
        </li>
      </ul>
      <div className="center actions">
        <Button outline onClick={onCancel}>close</Button>
      </div>
      {styles}
      <style jsx>{`
        .select-type {
          padding: 20px;
        }
        .center {
          text-align: center;
        }
        h1 {
          margin: 0;
        }
        .items {
          list-style: none;
          padding: 0;
          max-width: 360px;
          width: 100%;
          margin: 0 auto;
          padding-top: 30px;
        }
        .item {
        }
        .option {
          position: relative;
          display: block;
          width: 100%;
          padding: 20px;
          text-align: left;
          background: white;
          box-shadow: 0 3px 8px #efefef;
          color: #333;
          text-decoration: none;
          font-size: inherit;
          border: none;
        }
        .option:disabled {
          color: #999;
        }
        .option:hover {
        }
        .option:focus {
          z-index: 1;
        }
        .option .title {
          display: block;
        }
        .option .separator {
          display: none;
        }
        .option .description {
          display: block;
          padding-top: 5px;
        }
        .option .label {
          display: inline-block;
          font-size: 11px;
          font-style: normal;
          background: #5335B8;
          color: white;
          padding: 3px 8px;
          border-radius: 6px;
          position: absolute;
          top: 12px;
          right: 15px;
        }
        .option:disabled .label {
          background: #333;
        }
        .option .label.alert {
          background: #914146;
        }
        .actions {
          padding-top: 40px;
        }
      `}</style>
    </div>
  );
}
