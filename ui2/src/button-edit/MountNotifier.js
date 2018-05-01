import * as React from "react";

export default class MountNotifier extends React.Component {

    // Used to determine when a modal opens
    componentDidMount() {
        this.props.componentDidMount();
    }

    render() {
        return '';
    }
}
