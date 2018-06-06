import * as React from "react";
import './LedSegmentChooser.css';
import {Checkbox, Image, Modal} from "semantic-ui-react";

export default class LedSegmentChooser extends React.Component {

    render() {
        return <Modal trigger={this.props.trigger}
                      header='Select Light Segment'
                      content={<div className='scrolling content'>
                          <div className='description'>
                              {this.props.allSegments.map(seg => this.renderSegment(seg))}
                          </div>
                      </div>}
                      actions={[
                          {key: 'ok', content: 'OK', primary: true},
                          {key: 'cancel', content: 'Cancel'}
                      ]}
                      onActionClick={e => e.target.textContent === 'OK' ?
                          (this.props.onOk !== undefined && this.props.onOk()) :
                          this.props.onCancel()}>
        </Modal>
    }

    renderSegment(seg) {
        return <div className='ui image lsc-led-segment-list'
                    key={seg.name}>
            <Image label={seg.label}
                   src={"/segment-icons/" + encodeURIComponent(seg.name) + ".svg"}/>
            <div className='lsc-check-mark'>
                <Checkbox checked={this.props.checkedSegmentNames.includes(seg.name)}
                          onChange={() => this.props.toggleCheckedSegment(seg)}/>
            </div>
        </div>;
    }
}
