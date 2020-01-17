import React from 'react';
import ReactDOM from 'react-dom';

import NameForm from './NameForm.jsx';

class Main extends React.Component {

    constructor(props) {
        super(props);
        this.promises = {};
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
            const {type, msg} = JSON.parse(evt.data);

            this.setState({dataFromServer: type});
            console.log(type, msg);

            const {id} = msg;
            if (id in this.promises) {
                const {reject, resolve} = this.promises[id];
                delete this.promises[id];
                if (type === 'ERROR_REPORT') {
                    reject(msg.msg);
                } else {
                    msg.type = type;
                    resolve(msg);
                }
            } else {
                console.warn(`no promise for id ${id}`);
            }
        };

        this.ws.onclose = () => {
            console.log('disconnected')
            // automatically try to reconnect on connection loss
        };
    }

    sendPacket(type, msg) {
        const id = msg.id = Math.floor(Math.random() * 2147483648);

        const body = JSON.stringify({
            type: type,
            msg: msg,
        });
        this.ws.send(body);

        return new Promise((resolve,reject) => {
            if (this.promises[id]) {
                throw `Conflicting message id: ${id}`;
            }
            this.promises[id] = {
                resolve: resolve,
                reject: reject
            };
        });
    }

    async changeName(name) {
        await this.sendPacket('CHANGE_NAME', {
            new_name: name,
        });
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