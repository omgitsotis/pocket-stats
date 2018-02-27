import React, {Component} from 'react';

class Menu extends Component {
    constructor(props) {
        super(props);
        this.state = {
            currentPage: 'menu'
        };
    }

    onClick = (page) => {
        this.setState({currentPage: page});
    }

    render() {
        return (
            <div>
                <button>Initalise</button>
            </div>
        )
    }
}

export default Menu;
