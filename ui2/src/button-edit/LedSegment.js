import * as React from "react";
import {Image, Label, Segment} from "semantic-ui-react";
import './LedSegment.css';

export default class LedSegment extends React.Component {

    render() {
        let segment = this.props.segment;
        return <Segment attached color='blue' className='led-segment-container'>
            <Label content={segment.segment}
                   onRemove={() => this.props.onRemove(segment)}
                   color='blue' attached='top left'/>
            <Image className='led-segment-image'
                   size='tiny'
                   inline
                   src={"/segment-icons/" + encodeURIComponent(segment.segment) + ".svg"}/>
            <span>{segment.animation}<br/>{segment.params}</span>
        </Segment>;
    }
}
