import React, {Component} from 'react';
import Login from './login/Login.jsx';
import Socket from './socket.js';

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            authorised : false
        };
    }

    componentDidMount() {
        let ws = this.ws = new WebSocket('ws://localhost:4000')
        let socket = this.socket = new Socket(ws);
        socket.on('send auth', this.onAuth.bind(this));
        socket.on('subscribe auth', this.onRecievedAuth.bind(this));
    }

    onClick(e) {
        this.socket.emit('send auth');
    }

    onAuth(auth) {
        console.log(auth);
        window.open(auth.url, "myWindow", 'width=800,height=600');
    }

    onRecievedAuth() {
        this.setState({authorised: true});
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

export default App
