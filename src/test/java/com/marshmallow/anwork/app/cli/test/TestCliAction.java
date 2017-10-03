package com.marshmallow.anwork.app.cli.test;

import com.marshmallow.anwork.app.cli.CliAction;
import com.marshmallow.anwork.app.cli.CliFlags;

/**
 * This is a dummy {@link CliAction} class to be used in CLI tests.
 *
 * <p>
 * Created Sep 9, 2017
 * </p>
 *
 * @author Andrew
 */
public class TestCliAction implements CliAction {

  private boolean ran = false;
  private CliFlags flags = null;
  private String[] arguments = null;

  @Override
  public void run(CliFlags flags, String[] arguments) {
    ran = true;
    this.flags = flags;
    this.arguments = arguments;
  }

  public boolean getRan() {
    return ran;
  }

  public CliFlags getFlags() {
    return flags;
  }

  public String[] getArguments() {
    return arguments;
  }

}
