import React from 'react';
import ReactDOM from 'react-dom';

import NameForm from './NameForm.jsx';

class Main extends React.Component {

    constructor(props) {
        super(props);

        this.ws = new WebSocket(`ws://${location.host}/ws/`);
        this.state = {
            name: '?'
        };
    }

    componentDidMount() {
        this.ws.onopen = () => {
            console.log("connected");
        };

        this.ws.onmessage = evt => {
            // listen to data sent from the websocket server
            const message = JSON.parse(evt.data);
            this.setState({dataFromServer: message});
            console.log(message);
        };

        this.ws.onclose = () => {
            console.log('disconnected')
            // automatically try to reconnect on connection loss
        };
    }

    changeName(name) {
        const body = JSON.stringify({
            type: 'CHANGE_NAME',
            msg: {
                NewName: name,
            }
        });
        this.ws.send(body);
        console.log(body);
        this.setState({name: name});
    }

    render() {
        return (<div>
            <h>Hello {this.state.name}</h>
            <br/>
            <NameForm onSubmit={name => this.changeName(name)}/>
            <b>{JSON.stringify(this.state.dataFromServer)}</b>
        </div>);
    }
}

ReactDOM.render(
    <Main/>,
    document.getElementById('root')
);