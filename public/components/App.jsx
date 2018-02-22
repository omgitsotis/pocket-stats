import React, {Component} from 'react';
import Login from './login/Login.jsx';
import Socket from './socket.js';
import { instanceOf } from 'prop-types';
import { withCookies, Cookies } from 'react-cookie';

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            authorised : false,
        };
    }

    componentDidMount() {
        let ws = this.ws = new WebSocket('ws://localhost:4000')
        let socket = this.socket = new Socket(ws);
        socket.on('send auth', this.onAuth.bind(this));
        socket.on('subscribe auth', this.onRecievedAuth.bind(this));
        socket.on('data get', this.onDataGet.bind(this));
        socket.on('auth cached', this.onAuthCached.bind(this));

        const { cookies } = this.props;
        const accessToken = cookies.get('accessToken');
        if (typeof  accessToken !== "undefined") {
            console.log(accessToken);
            this.waitForSocketConnection(accessToken);
        }
    }

    waitForSocketConnection(accessToken) {
        setTimeout(
            () => {
                if (this.ws.readyState === 1) {
                    console.log("Connection is made")
                    this.socket.emit('auth cached', {token: accessToken});
                    return;
                } else {
                    console.log("wait for connection...")
                    this.waitForSocketConnection(accessToken);
                }
            }, 5);
    }

    onClick(e) {
        this.socket.emit('send auth');
    }

    onAuth(auth) {
        console.log(auth);
        window.open(auth.url, "myWindow", 'width=800,height=600');
    }

    onAuthCached() {
        this.setState({authorised: true});
    }

    onRecievedAuth(accessToken) {
        const { cookies } = this.props;
        const token = cookies.get('accessToken');
        if (typeof  token === "undefined") {
            cookies.set('accessToken', accessToken, { path: '/' });
        }

        this.setState({authorised: true});
        this.socket.emit('data get');
    }

    onDataGet(data) {
        console.log(data);
    }

    render() {
        return (
            <div className='app container'>
                <div>Text here</div>
                {this.state.authorised ?
                    (
                        <div>Authed</div>
                    ) : (
                    <Login onClick={this.onClick.bind(this)} />
                )}
            </div>
        )
    }
}

export default withCookies(App);
