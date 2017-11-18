package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Action;
import com.marshmallow.anwork.app.cli.ArgumentValues;

/**
 * This is an {@link Action} that simply prints out the version of the ANWORK app.
 *
 * <p>
 * Created Nov 18, 2017
 * </p>
 *
 * @author Andrew
 */
public class AnworkVersionAction implements Action {

  // This variable is updated by a gradle script in source! It reflects the "version" property
  // in the build.gradle script.
  private static final int VERSION = 1;

  @Override
  public void run(ArgumentValues flags, ArgumentValues values) {
    System.out.println("ANWORK Version = " + VERSION);
  }
}
