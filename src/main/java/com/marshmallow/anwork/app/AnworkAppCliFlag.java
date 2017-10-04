package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.CliArgumentType;

/**
 * These are the {@link CliFlag}'s that exist in the ANWORK CLI API. Each flag has a short flag
 * associated with it and maybe a {@link CliArgumentType} if it takes a parameter.
 *
 * <p>
 * Created Oct 4, 2017
 * </p>
 *
 * @author Andrew
 */
public enum AnworkAppCliFlag {
  CONTEXT("c", CliArgumentType.STRING),
  PERSISTENCE_ROOT("o", CliArgumentType.STRING),
  DONT_PERSIST("n"),
  DEBUG("d")
  ;

  private final String shortFlag;
  private final CliArgumentType parameterType;

  private AnworkAppCliFlag(String shortFlag, CliArgumentType parameterType) {
    this.shortFlag = shortFlag;
    this.parameterType = parameterType;
  }

  private AnworkAppCliFlag(String shortFlag) {
    // By default, flags with no parameters are translated to BOOLEAN values.
    this.shortFlag = shortFlag;
    this.parameterType = CliArgumentType.BOOLEAN;
  }

  public String getShortFlag() {
    return shortFlag;
  }

  public CliArgumentType getParameterType() {
    return parameterType;
  }
}
