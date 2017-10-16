package com.marshmallow.anwork.app.cli;

import java.util.ArrayList;
import java.util.List;

/**
 * This is a helper class for CLI parsing functionality.
 *
 * <p>
 * Created Sep 10, 2017
 * </p>
 *
 * @author Andrew
 */
class ParseContext {

  private Node activeNode;
  private List<String> parameters = new ArrayList<String>();
  private ArgumentValues flags = new ArgumentValues();

  Node getActiveNode() {
    return activeNode;
  }

  void setActiveNode(Node activeNode) {
    this.activeNode = activeNode;
  }

  String[] getParameters() {
    return parameters.toArray(new String[0]);
  }

  void addParameter(String parameter) {
    parameters.add(parameter);
  }

  ArgumentValues getFlags() {
    return flags;
  }
}
