package com.marshmallow.anwork.app.cli;

import java.util.ArrayList;
import java.util.List;

/**
 * This is a helper class for CLI parsing functionality.
 *
 * @author Andrew
 * @date Sep 10, 2017
 */
class CliParseContext {

  private CliNode activeNode;
  private List<String> parameters = new ArrayList<String>();

  CliNode getActiveNode() {
    return activeNode;
  }

  void setActiveNode(CliNode activeNode) {
    this.activeNode = activeNode;
  }

  String[] getParameters() {
    return parameters.toArray(new String[0]);
  }

  void addParameter(String parameter) {
    parameters.add(parameter);
  }
}
