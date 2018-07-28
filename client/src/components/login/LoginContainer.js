import React, {Component}  from 'react';
import Login from './Login.jsx';

class LoginContainer extends Component {
    render() {
        return <Login onClick={this.props.onClick} />;
    }
}

export default LoginContainer;
