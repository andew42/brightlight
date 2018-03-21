import * as React from "react";
import {Label, Segment} from "semantic-ui-react";

export default class LedSegment extends React.Component {

    render() {
        let segment = this.props.segment;
        return <Segment attached color='blue'>
            <Label content={segment.segment}
                   onRemove={() => this.props.onRemove(segment)}
                   color='blue' attached='top left'/>
            <div>{segment.animation}<br/>{segment.params}</div>
        </Segment>;
    }
}
