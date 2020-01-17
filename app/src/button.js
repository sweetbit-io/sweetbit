import React from 'react';
import classnames from 'classnames';
import Spinner from './spinner';

export default function({ children, onClick, submit, loading, outline, href, ...rest }) {
  const Tag = href ? 'a' : 'button';
  const type = href ? null : (submit ? 'submit' : 'button');

  return (
    <Tag
      onClick={onClick}
      className={classnames('submit', { loading, outline })}
      type={type}
      href={href}
      {...rest}
    >
      <span className="label">{children}</span>
      <span className="spinner">
        <Spinner />
      </span>
      <style jsx>{`
        * {
          box-sizing: border-box;
        }
        .submit {
          color: white;
          border: none;
          display: inline-flex;
          border-radius: 6px;
          background: #5335B8;
          border: 2px solid #5335B8;
          font-size: 16px;
          line-height: 18px;
          text-decoration: none;
          white-space: nowrap;
          font-weight: 300;
          overflow: hidden;
          transition: opacity 0.3s ease;
          padding: 12px 20px;
          position: relative;
        }
        .submit.outline {
          border: 2px solid #5335B8;
          background: transparent;
          color: #5335B8;
        }
        .submit .label {
          opacity: 1;
          transition: opacity 0.3s ease;
        }
        .submit .spinner {
          position: absolute;
          top: 50%;
          left: 50%;
          margin-left: -12px;
          margin-top: -12px;
          opacity: 0;
          transition: opacity 0.3s ease;
        }
        .submit.loading .label {
          opacity: 0;
        }
        .submit.loading .spinner {
          opacity: 1;
        }
      `}</style>
    </Tag>
  );
}
