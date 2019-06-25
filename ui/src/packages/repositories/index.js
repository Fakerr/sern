import React, { Component } from 'react';
import { Collapse, List, Button } from 'antd';

const { Panel } = Collapse;

const node = (
    <p style={{ color: '#1890ff', textDecoration: 'underline', marginBottom: '0' }}>My Repositories</p>
);

const customPanelStyle = {
  background: '#f7f7f7',
  color: 'green'
};

export default class RepoList extends Component {

  constructor(props) {
    super(props);
    this.state = {
      data: []
    };
  }

  componentDidMount() {
    fetch('/api/repos')
      .then(response => response.json())
      .then(data => this.setState({ data }));
  }

  render() {
    return (
      <div>
        <Collapse defaultActiveKey={['1']} style={{ borderRadius: '0' }}>
          <Panel header={node} key="1" style={customPanelStyle}>
	    <List
              bordered
	      size='small'
              dataSource={this.state.data}
              renderItem={item => (
		<List.Item>
		  <Button value={item} type="link" size={'large'} onClick={this.props.selectItem}>
		    {item}
		  </Button>
		</List.Item>
	      )}
	    />
          </Panel>
        </Collapse>
      </div>
    );
  }
}
