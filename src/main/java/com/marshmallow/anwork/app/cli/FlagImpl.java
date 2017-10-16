package com.marshmallow.anwork.app.cli;

/**
 * This is the default implementation of {@link MutableFlag}.
 *
 * <p>
 * Created Oct 10, 2017
 * </p>
 *
 * @author Andrew
 */
class FlagImpl implements MutableFlag {

  private String shortFlag;
  private String longFlag;
  private String description;
  private MutableArgument argument;

  public FlagImpl(String shortFlag) {
    this.shortFlag = shortFlag;
  }

  @Override
  public MutableFlag setShortFlag(String shortFlag) {
    this.shortFlag = shortFlag;
    return this;
  }

  @Override
  public String getShortFlag() {
    return shortFlag;
  }

  @Override
  public MutableFlag setLongFlag(String longFlag) {
    this.longFlag = longFlag;
    return this;
  }

  @Override
  public boolean hasLongFlag() {
    return longFlag != null;
  }

  @Override
  public String getLongFlag() {
    return longFlag;
  }

  @Override
  public MutableFlag setDescription(String description) {
    this.description = description;
    return this;
  }

  @Override
  public boolean hasDescription() {
    return description != null;
  }

  @Override
  public String getDescription() {
    return description;
  }

  @Override
  public <T> MutableArgument setArgument(String name, ArgumentType<T> type) {
    argument = new ArgumentImpl(name, type);
    return argument;
  }

  @Override
  public boolean hasArgument() {
    return argument != null;
  }

  @Override
  public Argument getArgument() {
    return argument;
  }

  @Override
  public int compareTo(Flag otherFlag) {
    return shortFlag.compareTo(otherFlag.getShortFlag());
  }

  @Override
  public String toString() {
    StringBuilder builder = new StringBuilder();
    builder.append(getClass().getName());
    builder.append(':').append(getShortFlag());
    if (hasLongFlag()) {
      builder.append(':').append(getLongFlag());
    }
    if (hasDescription()) {
      builder.append(':').append(getDescription());
    }
    if (hasArgument()) {
      builder.append(':').append(getArgument());
    }
    return builder.toString();
  }
}
