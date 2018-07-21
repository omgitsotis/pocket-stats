import React, {Component} from 'react';
import { withCookies, Cookies } from 'react-cookie';

import Socket from './socket.js';

import LoginContainer from './components/login/LoginContainer.js';
import DashboardContainer from './components/dashboard/DashboardContainer.jsx';

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            authorised : false,
            loaded: false,
            isDataError: false,
            updateComplete: true,
            lastUpdated: 0,
            username: "",
            totals: {},
            itemisedDate: {},
            itemisedTags: {},
        };
    }

    componentDidMount() {
        let ws = this.ws = new WebSocket('ws://localhost:4000')
        let socket = this.socket = new Socket(ws);

        socket.on('auth link', this.onAuthLink);
        socket.on('auth user', this.onAuthUser);
        socket.on('auth cached', this.onAuthCached);
        socket.on('data get', this.onDataGet);
        socket.on('data update', this.onDataUpdate);
        socket.on('error', this.onError);

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

    onClick = (e) => {
        this.socket.emit('send auth');
    }

    onAuthLink = (auth) => {
        console.log(auth);
        window.open(auth.url, "myWindow", 'width=800,height=600');
    }

    onAuthCached = (user) => {
        this.setState({
            authorised: true,
            username: user.username,
            lastUpdated: user.last_updated
        });

        const { cookies } = this.props;
        let token = cookies.get('accessToken');
        let userID = cookies.get('userID');

        this.socket.emit('data load', {
            token: token,
            id: parseInt(userID, 10),
        });
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

    onFetchDataClick = (start, end) => {
        const { cookies } = this.props;
        const userID = cookies.get('userID');
        this.socket.emit('data get', {
            start: start,
            end: end
        });
    }

    onAuthUser = (user) => {
        const { cookies } = this.props;
        let token = cookies.get('accessToken');
        let userID = cookies.get('userID');

        if (typeof token === "undefined") {
            cookies.set('accessToken', user.access_token, { path: '/' });
            cookies.set('userID', user.id, { path: '/' });

            token = user.access_token;
            userID = user.id;
        }

        //test
        this.setState({
            authorised: true,
            lastUpdated: user.last_updated
        });

        this.socket.emit('data load', {
            token: token,
            id: parseInt(userID, 10),
        });
    }

    onDataGet = (data) => {
        console.log(data);
        this.setState({
            loaded: true,
            updateComplete: true,
            totals: data.totals,
            itemisedDate: data.date_values,
            itemisedTags: data.tag_values
        });
    }

    onError = (err) => {
        console.error("error in", err.hookname, err.msg);
        switch (err.hookname) {
            case "data get":
                this.setState({
                    isDataError: true,
                    loaded: true,
                });
                break;
            default:

        }
    }

    onUpdateClick = () => {
        const { cookies } = this.props;
        const token = cookies.get('accessToken');
        const userID = cookies.get('userID');

        this.socket.emit('data update', {
            token: token,
            id: parseInt(userID, 10),
        });

        this.setState({
            updateComplete: false
        });
    }

    onDataUpdate = (data) => {
        console.log(data)
        this.setState({
            lastUpdated: data
        });

        const { cookies } = this.props;
        let token = cookies.get('accessToken');
        let userID = cookies.get('userID');

        this.socket.emit('data load', {
            token: token,
            id: parseInt(userID, 10),
        });
    }

    render() {
        let component;
        if (!this.state.authorised) {
            component = <LoginContainer onClick={this.onClick} />;
        } else {
            if(this.state.loaded === true) {
                component =
                    <DashboardContainer
                        totals={this.state.totals}
                        itemisedDate={this.state.itemisedDate}
                        itemisedTags={this.state.itemisedTags}
                        lastUpdated={this.state.lastUpdated}
                        onUpdateClick={this.onUpdateClick}
                        onFetchDataClick={this.onFetchDataClick}
                        updateComplete={this.state.updateComplete}
                        isDataError={this.state.isDataError}
                    />
            } else {
                component = <i className="fa fa-spinner fa-spin" style={{fontSize: "48px"}}></i>
            }
        }

        return (
            <div className='app container'>
                {component}
            </div>
        );
    }
}

export default withCookies(App);
