import React, {Component}  from 'react'
import './login.css';

class Login extends Component {
    render() {
        return (
            <div className='col-lg'>
                <h1>Pocket Stats</h1>
                <div className="btn-row">
                    <button className='btn btn-primary btn-lg'
                        onClick={this.props.onClick}>
                        Login
                    </button>
                </div>
            </div>
        )
    }
}

export default Login
