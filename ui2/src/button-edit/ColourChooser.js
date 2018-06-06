import * as React from "react";
import './ColourChooser.css';
import {Modal} from "semantic-ui-react";
import "../colour/Colour.js"
import VerticalSlider from "../buttons/VerticalSlider";
import MountNotifier from "./MountNotifier";
import Colour from "../colour/Colour";

export default class ColourChooser extends React.Component {
    // props: colour, trigger, onColourChanged
    render() {
        let colour = new Colour(this.props.colour);
        let saturatedColour = new Colour({h: colour.h, s: 100, v: 100});
        return <Modal trigger={this.props.trigger}
                      header='Choose Colour'
                      content={<div className='scrolling content'>
                          <MountNotifier componentDidMount={() => this.initialColour = colour}/>
                          <div className='cp-content-height'>
                              <ColourSwatch className='cp-colour-swatch' colour={colour.asColourString()}/>
                              <div className='cp-sliders'>
                                  <VerticalSlider className='cp-slider cp-slider-hue'
                                                  min={0}
                                                  max={360}
                                                  pos={colour.h}
                                                  onPosChange={p => this.props.onColourChanged(
                                                      new Colour({h: p, s: colour.s, v: colour.v}))}/>
                                  <VerticalSlider className='cp-slider cp-slider-saturation'
                                                  min={0}
                                                  max={100}
                                                  pos={colour.s}
                                                  onPosChange={p => this.props.onColourChanged(
                                                      new Colour({h: colour.h, s: p, v: colour.v}))}
                                                  sliderColour={saturatedColour.asColourString()}/>
                                  <VerticalSlider className='cp-slider cp-slider-value'
                                                  min={0}
                                                  max={100}
                                                  pos={colour.v}
                                                  onPosChange={p => this.props.onColourChanged(
                                                      new Colour({h: colour.h, s: colour.s, v: p}))}
                                                  sliderColour={saturatedColour.asColourString()}/>
                              </div>
                          </div>
                      </div>}
                      actions={[
                          {key: 'ok', content: 'OK', primary: true},
                          {key: 'cancel', content: 'Cancel'}
                      ]}
                      onActionClick={e => e.target.textContent === 'Cancel' ?
                          this.props.onColourChanged(this.initialColour) :
                          null}/>
    }
}

function ColourSwatch(props) {
    return <div className='cp-colour-swatch' style={{backgroundColor: props.colour}}/>
}
