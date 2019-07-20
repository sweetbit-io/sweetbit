import React, { Component } from 'react'

export default class Return extends Component {
  timeout = null

  constructor(props) {
    super()
    this.state = {
      seconds: props.seconds,
    }
  }

  componentDidMount() {
    if (this.props.seconds !== null) {
      this.timeout = setTimeout(() => this.tick(), 1000)
    }
  }

  componentDidUpdate(prevProps) {
    if (this.timeout) {
      clearTimeout(this.timeout)
    }

    if (prevProps.seconds !== this.props.seconds) {
      this.setState({
        seconds: this.props.seconds,
      })
    }

    if (this.props.seconds !== null) {
      this.timeout = setTimeout(() => this.tick(), 1000)
    }
  }

  componentWillUnmount() {
    if (this.timeout) {
      clearTimeout(this.timeout)
    }
  }

  tick() {
    if (this.state.seconds > 0) {
      this.setState({
        seconds: this.state.seconds - 1,
      })
    } else {
      this.props.onReturn()
    }
  }

  render() {
    return this.props.children(this.state.seconds)
  }
}
