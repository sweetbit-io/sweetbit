import React, { useCallback } from 'react';

export default function Toggle({
  name,
  checked,
  value,
  disabled,
  onChange,
}) {
  const change = useCallback((e) => {
    onChange(e.target.checked);
  }, [onChange]);

  return (
    <label className="toggle">
      <input
        type="checkbox"
        name={name}
        checked={checked}
        value={value}
        disabled={disabled}
        onChange={change}
      />
      <div className="slider" />
      <style jsx>{`
        .toggle {
          position: relative;
          display: flex;
          align-items: center;
          margin: 0;
        }
        input {
          display: none;
        }
        .slider {
          display: block;
          margin-left: 10px;
          position: relative;
          cursor: pointer;
          width: 60px;
          height: 34px;
          background-color: #ccc;
          transition: .4s;
          border-radius: 34px;
          border: solid 2px #ccc;
        }
        .slider:before {
          position: absolute;
          content: '';
          background-color: #fff;
          transition: .4s;
          border: 0;
          border-radius: 50%;
          height: 28px;
          width: 28px;
          top: 1px;
          left: 1px;
        }
        .slider:after {
          position: absolute;
          content: url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 12 12"><path d="M11.54 2.828l-6.064 8.72c-.195.28-.489.443-.808.452-.335 0-.626-.148-.823-.408L1.246 8.22a1.216 1.216 0 0 1-.242-.838c.023-.306.153-.585.366-.783.412-.383 1.13-.316 1.477.135l1.757 2.281 5.27-7.582c.329-.471 1.039-.577 1.469-.217.223.186.365.457.402.761.037.306-.036.607-.205.852z" fill="#2E4DB9"/></svg>');
          display: flex;
          align-items: center;
          justify-content: center;
          height: 28px;
          width: 28px;
          top: 1px;
          left: 1px;
          opacity: 0;
          transition: .4s;
        }
        input:checked + .slider {
          background-color: #3032c1;
          border-color: #3032c1;
        }
        input:focus + .slider {
          box-shadow: 0 0 1px #3032c1;
        }
        input:checked + .slider:before {
          transform: translateX(25px);
          box-shadow: 0px 0px 0px 8px rgba(48,50,193,.2);
        }
        input:checked + .slider:after {
          transform: translateX(25px);
          opacity: 1;
        }
      `}</style>
    </label>
  );
}
