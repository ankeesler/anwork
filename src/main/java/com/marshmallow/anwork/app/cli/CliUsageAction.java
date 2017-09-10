package com.marshmallow.anwork.app.cli;

/**
 * This is a command action that simply prints out the usage of a CLI node.
 *
 * @author Andrew
 * @date Sep 10, 2017
 */
public class CliUsageAction implements CliAction {

  private CliNode node;

  public CliUsageAction(CliNode node) {
    this.node = node;
  }

  @Override
  public void run(String[] arguments) {
    System.out.println(node.getUsage());
  }
}
