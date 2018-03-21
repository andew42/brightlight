import * as React from "react";
import './LedSegmentList.css';
import {Button, Icon, Image, Modal} from "semantic-ui-react";

export default class LedSegmentList extends React.Component {

    constructor(props) {
        super(props);

        this.state = {isOpen: true};
    }

    renderSegment(seg) {
        return <div className='ui image led-segment-list'>
            <Image label={seg} src={"/segment-icons/" + encodeURIComponent(seg) + ".svg"}/>
            <div className='check-mark'>
                <Icon name='checkmark'/>
            </div>
        </div>;
    }

    render() {
        let segments = ['All', 'All Ceiling', 'All Wall', 'Bathroom', 'Bathroom Ceiling', 'Bathroom Wall',
            'Bedroom', 'Bedroom Ceiling', 'Bedroom Wall', 'Chest', 'Chest Ceiling', 'Chest Wall',
            'Curtains', 'Door', 'Dressing', 'Dressing Ceiling', 'Dressing Wall', 'Light Switch'];
        return <Modal trigger={this.props.trigger}
                      open={this.state.isOpen}
        >
            <Modal.Header>Select Light Segment</Modal.Header>
            <Modal.Content scrolling>
                <Modal.Description>
                    {segments.map(seg => this.renderSegment(seg))}
                </Modal.Description>
            </Modal.Content>
            <Modal.Actions>
                <Button primary content='OK'/>
                <Button content='Cancel' onClick={() => this.setState({isOpen: false})}/>
            </Modal.Actions>
        </Modal>
    }
}
