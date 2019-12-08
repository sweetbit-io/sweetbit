import React from 'react';
import css from 'styled-jsx/css';
import { Link, Switch, Route } from 'react-router-dom';
import Home from './home';
import Nodes from './nodes';
import Network from './network';
import Updates from './updates';
import { ReactComponent as HomeIcon } from './home.svg';
import { ReactComponent as NodeIcon } from './nodes.svg';
import { ReactComponent as NetworkIcon } from './network.svg';
import { ReactComponent as UpdateIcon } from './updates.svg';

const { className, styles } = css.resolve`
  a {
    display: flex;
    align-items: center;
    color: green;
  }
  svg {
    width: 32px;
    height: 32px;
  }
`

function Dispenser() {
  return (
    <div className="dispenser">
      <div className="menu">
        <ul>
          <li>
            <Link className={className} to="/">
              <HomeIcon className={className} />
              <span className="label">Candy Dispenser</span>
            </Link>
          </li>
          <li>
            <Link className={className} to="/nodes">
              <NodeIcon className={className} />
              <span className="label">Nodes</span>
            </Link>
          </li>
          <li>
            <Link className={className} to="/network">
              <NetworkIcon className={className} />
              <span className="label">Network</span>
            </Link>
          </li>
          <li>
            <Link className={className} to="/updates">
              <UpdateIcon className={className} />
              <span className="label">Updates</span>
            </Link>
          </li>
        </ul>
      </div>
      <div className="main">
        <Switch>
          <Route path="/" exact component={Home} />
          <Route path="/nodes" component={Nodes} />
          <Route path="/network" component={Network} />
          <Route path="/updates" component={Updates} />
        </Switch>
        <div className="">
        </div>
      </div>
      {styles}
      <style jsx>{`
        .dispenser {
          padding: 0 20px;
          margin: 0 auto;
          max-width: 768px;
          width: 100%;
          display: flex;
        }

        .menu {
          flex: 200px 0 0;
        }

        .menu ul {
          list-style: none;
          padding: 0;
          margin: 0;
        }

        .main {
          flex: auto 1 1;
        }
      `}</style>
    </div>
  );
}

export default Dispenser;
