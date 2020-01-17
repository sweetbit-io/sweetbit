import React, { useRef, useCallback } from 'react';
import css from 'styled-jsx/css';
import Button from './button';
import { formClassName, formStyles } from './form';
import { ReactComponent as UnlockImage } from './unlock.svg';

const { className, styles } = css.resolve`
  .image {
    width: auto;
    height: 120px;
  }
`;

export default function UnlockNode({
  onSetPassword,
  onCancel,
}) {
  const passwordEl = useRef(null);
  const submit = useCallback((e) => {
    e.preventDefault();
    onSetPassword(passwordEl.current.value);
  }, [passwordEl, onSetPassword]);

  return (
    <div className="unlock">
      <p className="center">
        <UnlockImage className={`${className} image`} />
      </p>
      <h1 className="center">Unlock your node</h1>
      <form onSubmit={submit}>
        <div className="name">
          <div className={`${formClassName} group centered`}>
            <input
              className={formClassName}
              ref={passwordEl}
              placeholder="Password"
              required
              name="password"
              type="password"
              autoComplete="password"
            />
            <label className={formClassName} htmlFor="password">Password</label>
          </div>
        </div>
        <div className="actions center">
          <Button submit type="submit">unlock</Button>
          <span> </span>
          <Button outline onClick={onCancel}>cancel</Button>
        </div>
      </form>
      {styles}
      <style jsx>{`
        .unlock {
          padding: 20px;
        }

        .center {
          text-align: center;
        }

        .actions {
          padding-top: 40px;
        }

        h1 {
          margin: 0;
        }
      `}</style>
      {formStyles}
    </div>
  );
}
