import React, { Component } from 'react';

class Monitor extends Component {
    constructor(props) {
        super(props)
        this.state = {
            drone : null,
            server : null,
        }
    }

    componentDidMount() {
        this.time = setInterval(
          () => this.tick(),
          1000
        );
    }

    componentWillUnmount() {
        clearInterval(this.timer);
    }

    tick() {
        fetch("http://3.91.1.79:10000/state")
        .then( res => res.json())
        .then( 
            (result) => {
                console.log(result)
                this.setState({
                    drone : {
                        dji : result.drone.dji,
                        saturatr : result.drone.saturatr
                    },
                    server : {
                        saturatr : result.server.saturatr
                    }
                });
            },
            (error) => {
                console.error("Problem loading state")
            }
        )
    }

    render() {
        return (
            <div>
                <div>
                    <h2>Drone</h2>
                    <Dji status={this.state.drone.dji}></Dji>
                    <Saturatr status={this.state.drone.saturatr}></Saturatr>
                </div>
                <div>
                    <h2>Server</h2>
                    <Saturatr status={this.state.server.saturatr}></Saturatr>
                </div>
            </div>
        )
    }
}

export default Monitor

class Dji extends Component {
    render() {
        return (
            <div>
                <h3>Dji</h3>
                <p>battery: {this.props.status.battery} </p>
                <p>altitude: {this.props.status.altitude} </p>
            </div>
        )
    }
}

class Saturatr extends Component {
    render() {
        return (
            <div>
                <h3>Saturatr</h3>
                <div>
                    <h4>Acker</h4>
                    <p>sent: {this.props.status.acker.sent}</p>
                    <p>received: {this.props.status.acker.received}</p>
                </div>
                <div>
                    <h4>Saturatr</h4>
                    <p>sent: {this.props.status.saturatr.sent}</p>
                    <p>received: {this.props.status.saturatr.received}</p>
                </div>
            </div>
        )
    }
}
