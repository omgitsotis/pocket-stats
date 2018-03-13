import React, {Component} from 'react';
import { withCookies, Cookies } from 'react-cookie';

import Socket from './socket.js';

import LoginContainer from './components/login/LoginContainer.js';
import DashboardContainer from './components/dashboard/Dashboard.jsx';

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            authorised : false,
            initState: '',
            startDate: 0,
            endDate: 0,
            loaded: false,
            statList: {}
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
        socket.on('data update', this.onDataUpdate);
        socket.on('data load', this.onDataLoad);

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

        this.socket.emit('data load', {
            token: token,
            id: parseInt(userID, 10),
        });
    }

    onDataGet(data) {
        console.log(data);
    }

    onError(err) {
        console.log("there was an error:", err);
    }

    onUpdateClick = () => {
        const { cookies } = this.props;
        const token = cookies.get('accessToken');
        const userID = cookies.get('userID');

        this.socket.emit('data update', {
            token: token,
            id: parseInt(userID, 10),
        });
    }

    onDataUpdate = (data) => {
        console.log(data)
    }

    onDataLoad = (data) => {
        console.log("On data load", data);
        this.setState({
            loaded: true,
            startDate: data.start_date,
            endDate: data.end_date,
            statList: {}
        });
    }

    render() {
        let component;
        if (!this.state.authorised) {
            component = <LoginContainer onClick={this.onClick.bind(this)} />;
        } else {
            if(this.state.loaded === true) {
                component = <DashboardContainer />
            } else {
                component = <i className="fa fa-spinner fa-spin" style={{fontSize: "48px"}}></i>
            }
        }

        return (
            <div className='app container'>
                <DashboardContainer />
            </div>
        );
    }
}

export default withCookies(App);
