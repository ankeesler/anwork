package com.marshmallow.anwork.app.cli.test;

import com.marshmallow.anwork.app.cli.CliAction;

/**
 * This is a dummy {@link CliAction} class to be used in CLI tests.
 *
 * @author Andrew
 * @date Sep 9, 2017
 */
public class TestCliAction implements CliAction {

  private boolean ran = false;
  private String[] arguments = null;

  @Override
  public void run(String[] arguments) {
    ran = true;
    this.arguments = arguments;
  }

  public boolean getRan() {
    return ran;
  }

  public String[] getArguments() {
    return arguments;
  }

}
