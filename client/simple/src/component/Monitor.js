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
		result.drone.dji.altitude = result.drone.dji.gps[2]
                this.setState({
                    drone : {
                        dji : result.drone.dji,
                        saturatr : result.drone.saturatr,
			signal: result.drone.signal, 
			upload : result.drone.upload,
			download : result.drone.download
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
                    { this.state.drone && 
                        <div>
                            <Dji status={this.state.drone.dji}></Dji>
                            <Saturatr status={this.state.drone.saturatr}></Saturatr>
			    <Signal status={this.state.drone.signal}></Signal>
			    <p>upload: {this.state.drone.upload} Mbps</p>
			    <p>download: {this.state.drone.download} Mbps</p>
                        </div>
                    }
                </div>
                <div>
                    <h2>Server</h2>
                    { this.state.server && 
                        <Saturatr status={this.state.server.saturatr}></Saturatr>
                    }
                </div>
            </div>
        )
    }
}

export default Monitor

class Signal extends Component {
    render() {
	return (
	    <div>
		<h3>Hilink</h3>
		<p>rsrp: {this.props.status.rsrp} </p>
		<p>rsrq: {this.props.status.rsrq} </p>
		<p>rssi: {this.props.status.rssi} </p>
		<p>sinr: {this.props.status.sinr} </p>
		<p>cell_id: {this.props.status.cell_id} </p>
	    </div>

	)
    }

}

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
