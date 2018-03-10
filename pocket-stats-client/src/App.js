import React, {Component} from 'react';
import Login from './components/login/Login.jsx';
import Socket from './socket.js';
import { withCookies, Cookies } from 'react-cookie';
import Menu from './components/menu/Menu.jsx';

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            authorised : false,
            initState: ''
        };
    }

    componentDidMount() {
        let ws = this.ws = new WebSocket('ws://localhost:4000')
        let socket = this.socket = new Socket(ws);
        socket.on('send auth', this.onAuth.bind(this));
        socket.on('subscribe auth', this.onRecievedAuth.bind(this));
        socket.on('data get', this.onDataGet.bind(this));
        socket.on('auth cached', this.onAuthCached.bind(this));
        socket.on('error', this.onError.bind(this));

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

    onInitClick = (date) => {
        const { cookies } = this.props;
        const token = cookies.get('accessToken');
        const userID = cookies.get('userID');

        console.log(date);
        this.setState({initState: 'started'});

        this.socket.emit('data init', {
            token: token,
            id: parseInt(userID, 10),
            date: date
        });
    }

    onFetchDataClick = () => {
        const { cookies } = this.props;
        const userID = cookies.get('userID');
        this.socket.emit('data get', {
            // id: parseInt(userID, 10),
            start: 1519084800,
            end: 1519776000
        });
    }

    onRecievedAuth(user) {
        const { cookies } = this.props;
        let token = cookies.get('accessToken');
        let userID = cookies.get('userID');
        if (typeof token === "undefined") {
            cookies.set('accessToken', user.access_token, { path: '/' });
            cookies.set('userID', user.id, { path: '/' });

            token = user.access_token;
            userID = user.id;
        }

        this.setState({authorised: true});
    }

    onDataGet(data) {
        console.log(data);
    }

    onError(err) {
        console.log("there was an error:", err);
    }

    render() {
        console.log(this.state.initState);
        return (
            <div className='app container'>
                {this.state.authorised ?
                    (
                        <Menu
                            onInitClick={ this.onInitClick }
                            initState={this.state.initState}
                            onFetchDataClick={this.onFetchDataClick}/>
                    ) : (
                    <Login onClick={this.onClick.bind(this)} />
                )}
            </div>
        )
    }
}

export default withCookies(App);
