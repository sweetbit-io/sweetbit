import React, { useRef, useCallback } from 'react';
import Button from '../button';

export default function RemoteLnd({
  onCancel,
  onAdd,
}) {
  const nameEl = useRef(null);
  const typeEl = useRef(null);
  const uriEl = useRef(null);
  const certEl = useRef(null);
  const macaroonEl = useRef(null);
  const submit = useCallback((e) => {
    e.preventDefault();
    onAdd({
      name: nameEl.current.value,
      type: typeEl.current.value,
      uri: uriEl.current.value,
      cert: certEl.current.innerText,
      macaroon: macaroonEl.current.value,
    });
  }, [nameEl, typeEl, uriEl, certEl, macaroonEl, onAdd]);

  return (
    <div>
      <h1>Add Lightning node</h1>
      <form onSubmit={submit}>
        <div className="name">
          <div className="group">
            <input
              ref={nameEl}
              placeholder="Name"
              required
              name="name"
              type="text"
              autoComplete="name"
            />
            <label htmlFor="name">Name</label>
          </div>
        </div>
        <div className="type">
          <div className="group">
            <select required ref={typeEl} defaultValue="">
              <option value="" disabled></option>
              <option value="remote-lnd">Remote LND node</option>
              <option value="local">Local LND node</option>
            </select>
            <label htmlFor="username">Type</label>
          </div>
        </div>
        <div className="uri">
          <div className="group">
            <input
              ref={uriEl}
              placeholder="URI"
              required
              name="uri"
              type="text"
              autoComplete="uri"
            />
            <label htmlFor="uri">URI</label>
          </div>
        </div>
        <div className="cert">
          <div className="group">
            <div contentEditable="true" ref={certEl}></div>
            <label htmlFor="cert">Certificate</label>
          </div>
        </div>
        <div className="macaroon">
          <div className="group">
            <input
              ref={macaroonEl}
              placeholder="Macaroon"
              required
              name="macaroon"
              type="text"
              autoComplete="macaroon"
            />
            <label htmlFor="macaroon">Macaroon</label>
          </div>
        </div>
        <div className="actions">
          <div className="action">
            <Button submit type="submit">add</Button>
          </div>
          <div className="action">
            <Button type="button" onClick={onCancel} outline>cancel</Button>
          </div>
        </div>
      </form>
      <style jsx>{`
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

        .group {
          position: relative;
          margin-top: 45px;
        }
        .group.centered input,
        .group.centered [contenteditable],
        .group.centered select {
          text-align: center;
        }
        .group input,
        .group textarea,
        .group [contenteditable],
        .group select {
          font-size: 18px;
          padding: 10px 10px 10px 5px;
          display: block;
          width: 100%;
          border: none;
          border-bottom: 3px solid #ccc;
          background: transparent;
          border-radius: 0;
          -webkit-appearance: none;
          -moz-appearance: none;
          appearance: none;
        }
        .group input:focus,
        .group textarea:focus,
        .group [contenteditable]:focus,
        .group select:focus {
          outline: none;
        }
        .group input::placeholder,
        .group textarea::placeholder,
        .group select::placeholder {
          color: transparent;
        }
        .group label {
          color: #999;
          font-size: 18px;
          font-weight: normal;
          position: absolute;
          pointer-events: none;
          left: 5px;
          top: 10px;
          transition: 0.2s ease all;
        }
        .group.centered label {
          left: 50%;
          transform: translateX(-50%);
        }
        .group input:focus ~ label,
        .group textarea:focus ~ label,
        .group [contenteditable]:focus ~ label,
        .group select:focus ~ label,
        .group input:not(:placeholder-shown) ~ label,
        .group textarea:not(:placeholder-shown) ~ label,
        .group input:valid ~ label,
        .group textarea:valid ~ label,
        .group [contenteditable]:not(:empty) ~ label,
        .group select:valid ~ label {
          top: -20px;
          font-size: 14px;
          color: #5264AE;
        }
      `}</style>
    </div>
  );
}
