import React, {Component}  from 'react'

class Login extends Component {
    render() {
        return (
            <button onClick={this.props.onClick}>Login</button>
        )
    }
}

export default Login
