import React, { Component } from 'react'
import Head from 'next/head'
import AnimateHeight from 'react-animate-height'
import QRCode from 'qrcode.react'
import classnames from 'classnames'
import Router from 'next/router'
import queryString from 'query-string'
import Candy from '../components/candy'
import Check from '../components/check'
import Cross from '../components/cross'
import Button from '../components/button'
import Return from '../components/return'

export default class IndexPage extends Component {
  state = {
    unavailable: false,
    notFound: false,
    invoice: null,
  }

  static getInitialProps({ query }) {
    return {
      apiBaseUrl: process.env.API_BASE_URL,
    }
  }

  async componentDidMount() {
    const { r_hash: rHash } = queryString.parse(location.search)

    await this.fetch(rHash, {
      apiBaseUrl: this.props.apiBaseUrl || `${window.location.origin}/api`,
    })
  }

  async fetch(rHash, { apiBaseUrl }) {
    let res;
    let invoice = null;

    if (rHash) {
      res = await fetch(`${apiBaseUrl}/invoices/${rHash}`, { method: 'GET' })
    } else {
      res = await fetch(`${apiBaseUrl}/invoices`, { method: 'POST' })
    }

    if (res.status === 503) {
      this.setState({
        unavailable: true,
        notFound: false,
        invoice: null,
      })
      return
    } else if (res.status === 404) {
      this.setState({
        unavailable: false,
        notFound: true,
        invoice: null,
      })
      return
    } else {
      invoice = await res.json()
      this.setState({
        unavailable: false,
        notFound: false,
        invoice,
      })
    }

    const url = `/?r_hash=${invoice.r_hash}`
    Router.push(url, url, { shallow: true })

    if (invoice && !invoice.settled) {
      const statusSocket = new WebSocket(`${apiBaseUrl.replace('http', 'ws')}/invoices/${invoice.r_hash}/status`);

      statusSocket.onmessage = ({ data }) => {
        const { settled } = JSON.parse(data) || {}

        this.setState((state) => ({
          invoice: state.invoice !== null ? {
            ...state.invoice,
            settled,
          } : null,
        }))

        if (settled) {
          // close after invoice was settled
          statusSocket.close()
        }
      };
    }
  }

  handleRetry = async (e) =>  {
    e.preventDefault();

    this.setState({
      unavailable: false,
      notFound: false,
      invoice: null,
    })

    await this.fetch(null, {
      apiBaseUrl: this.props.apiBaseUrl || `${window.location.origin}/api`,
    })
  }

  handleGenerateNewInvoice = async (e) =>  {
    e.preventDefault();
    await this.fetch(null, {
      apiBaseUrl: this.props.apiBaseUrl || `${window.location.origin}/api`,
    })
  }

  handleReturnFromPaidInvoice = async () => {
    await this.fetch(null, {
      apiBaseUrl: this.props.apiBaseUrl || `${window.location.origin}/api`,
    })
  }

  render() {
    return (
      <div className="section">
        <Head>
          <meta charSet="utf-8" />
          <meta name="viewport" content="initial-scale=1.0, width=device-width, maximum-scale=1.0, user-scalable=no" />
          <meta name="apple-mobile-web-app-capable" content="yes" />
          <meta name="mobile-web-app-capable" content="yes" />
          <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
          <meta name="format-detection" content="telephone=no" />
          <meta name="apple-mobile-web-app-title" content="Sweet ⚡️" />
          <title>Candy Dispenser</title>
          <meta name="description" content="Pay for your candy with Bitcoin over Lightning" />
          <meta name="MobileOptimized" content="320" />
          <meta name="theme-color" content="#ffffff" />
          <link rel="icon" href="/favicon.ico" type="image/x-icon" />
          <meta name="twitter:title" content="Lightning Candy Dispenser" />
          <meta name="twitter:description" content="Pay for your candy with Bitcoin over Lightning" />
          {/* <meta name="twitter:image" content={imageUrl} /> */}
          {/* <meta name="twitter:card" content="summary_large_image" /> */}
          <meta name="twitter:card" content="summary" />
          <meta property="og:title" content="Lightning Candy Dispenser" />
          <meta property="og:description" content="Pay for your candy with Bitcoin over Lightning" />
          {/* <meta name="og:image" content={imageUrl} /> */}
        </Head>
        <div className="title">
          Candy Dispenser
        </div>
        <div className="description">
          Please use the invoice below in order to dispense your candy.
        </div>
        <div className="qr">
          <div className={classnames('code', { show: this.state.unavailable })}>
            <Cross />
            <p>node unavailable</p>
            <Button onClick={this.handleRetry}>
              retry
            </Button>
          </div>
          <div className={classnames('code', { show: this.state.notFound })}>
            <Cross />
            <p>not found.</p>
            <Button onClick={this.handleGenerateNewInvoice}>
              generate new
            </Button>
          </div>
          <div className={classnames('code', 'loading', { show: !this.state.unavailable && !this.state.invoice })}>
            loading
          </div>
          <div className={classnames('code', 'loading', { show: this.state.invoice && this.state.invoice.settled })}>
            <Check />
            <Return
              seconds={this.state.invoice && this.state.invoice.settled ? 5 : null}
              onReturn={this.handleReturnFromPaidInvoice}
            >
              {(secondsLeft) => (
                <Button secondary onClick={this.handleReturnFromPaidInvoice}>
                  {secondsLeft > 0 ? (
                    <span>Back to start in {secondsLeft}s</span>
                  ) : (
                    <span>Back to start...</span>
                  )}
                </Button>
              )}
            </Return>
          </div>
          <a
            className={classnames('code', { show: this.state.invoice && !this.state.invoice.settled })}
            href={`lightning:${this.state.invoice && this.state.invoice.payment_request}`}
          >
            <QRCode
              style={{ width: '100%', height: '100%', display: 'block' }}
              size={256}
              renderAs="svg"
              value={this.state.invoice && this.state.invoice.payment_request || ''}
            />
          </a>
        </div>
        <div className={classnames('payreq', { show: this.state.invoice && !this.state.invoice.settled })}>
          <div className="invoice">Your invoice ⚡️</div>
          <pre>
            <code>
              {this.state.invoice && this.state.invoice.payment_request}
            </code>
          </pre>
        </div>
        <style jsx>{`
          * {
            box-sizing: border-box;
            font-family: sans-serif;
          }

          .section {
            background: white;
            box-shadow: 0 0 99px rgba(0,0,0,0.3);
            padding: 35px 15px;
            border-radius: 6px;
            margin: 0 auto;
            max-width: 480px;
          }

          @media (min-width: 768px) {
            .section {
              margin: 100px auto;
            }
          }

          .title {
            position: relative;
            font-size: 34px;
            font-weight: 100;
            text-align: center;
            padding-top: 0;
          }

          .back {
            position: absolute;
            top: 0;
            left: 0;
            display: block;
            background: transparent;
            border: none;
            font: inherit;
            margin: 0;
            height: 100%;
            padding: 20px 10px 0;
            color: #333;
            cursor: pointer;
          }

          .back svg {
            display: block;
            height: 100%;
            width: auto;
          }

          .description {
            font-size: 18px;
            font-weight: 100;
            text-align: center;
            padding-top: 20px;
            color: #333;
          }

          .candy {
            text-align: center;
            padding-top: 26px;
          }

          .amount {
            display: flex;
            padding-top: 40px;
            justify-content: center;
          }

          .less, .more {
            border: none;
            font-size: inherit;
            background: none;
            margin: 0;
            padding: 0;
            color: inherit;
            flex: 0 0 74px;
            width: 74px;
            height: 74px;
            color: green;
          }

          .less:disabled, .more:disabled {
            color: #666;
            cursor: not-allowed;
          }

          .less span, .more span {
            display: none;
          }

          .value {
            flex: 0 auto;
            padding: 0 40px;
            display: flex;
            flex-direction: column;
            justify-content: center;
            overflow: hidden;
          }

          .sat {
            font-size: 18px;
          }

          .usd {
            font-size: 32px;
          }

          .sat, .usd {
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            text-align: center;
          }

          .action {
            padding-top: 40px;
            text-align: center;
          }

          .qr {
            position: relative;
            text-align: center;
            padding-top: 40px;
          }

          .code {
            display: inline-block;
            width: 100%;
            max-width: 350px;
            height: auto;
            padding: 25px;
            box-shadow: 0 0 40px rgba(0,0,0,0.3);
            opacity: 0;
            pointer-events: none;
            transition: box-shadow 0.3s ease, opacity 0.3s ease;
          }

          .code.show {
            opacity: 1;
            pointer-events: all;
          }

          .code:hover {
            box-shadow: 0 0 60px rgba(0,0,0,0.3);
          }

          .code.loading {
            position: absolute;
            top: 40px;
            opacity: 0;
            transition: opacity 0.3s ease;
          }

          .code.loading.show {
            opacity: 1;
          }

          .code.loading:after {
            /* Make it a responsive square */
            content: '';
            display: block;
            padding-bottom: 100%;
          }

          .payreq {
            padding-top: 20px;
            opacity: 0;
            pointer-events: none;
          }

          .payreq.show {
            opacity: 1;
            pointer-events: all;
          }

          .invoice {
            margin: 0 auto;
            max-width: 350px;
            padding: 15px;
            color: #666;
            text-align: center;
          }

          pre {
            font-size: 20px;
            word-wrap: break-word;
            white-space: normal;
            padding: 25px;
            cursor: copy;
            transition: background-color .3s ease;
            margin: 0 auto;
            max-width: 350px;
            box-shadow: 0 0 40px rgba(0,0,0,0.3);
          }
        `}</style>
      </div>
    )
  }
}
