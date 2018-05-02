import * as React from "react";
import './LedSegmentChooser.css';
import {Checkbox, Image, Modal} from "semantic-ui-react";
import PropTypes from 'prop-types';

export default class LedSegmentChooser extends React.Component {

    static propTypes = {
        // segment names that are checked
        checkedSegmentNames: PropTypes.array.isRequired,

        // button used to open chooser
        trigger: PropTypes.object.isRequired,

        // list of segment names to choose from
        allSegmentNames: PropTypes.array.isRequired,

        // function to call when a segment is toggled
        toggleCheckedSegment: PropTypes.func.isRequired,

        // function to call when user hits OK
        onOk: PropTypes.func.isRequired,

        // function to call when user hits cancel
        onCancel: PropTypes.func.isRequired
    };

    renderSegment(seg) {
        return <div className='ui image led-segment-list'
                    key={seg.segment}>
            <Image label={seg.label}
                   src={"/segment-icons/" + encodeURIComponent(seg.segment) + ".svg"}/>
            <div className='check-mark'>
                <Checkbox checked={this.props.checkedSegmentNames.includes(seg.segment)}
                          onChange={() => this.props.toggleCheckedSegment(seg.segment)}/>
            </div>
        </div>;
    }

    render() {
        return <Modal trigger={this.props.trigger}
                      header='Select Light Segment'
                      content={<div className='scrolling content'>
                          <div className='description'>
                              {this.props.allSegmentNames.map(seg => this.renderSegment(seg))}
                          </div>
                      </div>}
                      actions={[
                          {key: 'ok', content: 'OK', primary: true},
                          {key: 'cancel', content: 'Cancel'}
                      ]}
                      onActionClick={e => e.target.textContent === 'OK' ? this.props.onOk() : this.props.onCancel()}>
        </Modal>
    }
}
