package com.marshmallow.anwork.app.cli;

/**
 * This is a command action that simply prints out the usage of a CLI node.
 *
 * <p>
 * Created Sep 10, 2017
 * </p>
 *
 * @author Andrew
 */
public class CliUsageAction implements CliAction {

  private CliNodeImpl node;

  public CliUsageAction(CliNodeImpl node) {
    this.node = node;
  }

  @Override
  public void run(CliFlags flags, String[] arguments) {
    System.out.println(node.getUsage());
  }
}
