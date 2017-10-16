package com.marshmallow.anwork.app.cli;

/**
 * This is the default implementation of {@link MutableArgument}.
 *
 * <p>
 * Created Oct 10, 2017
 * </p>
 *
 * @author Andrew
 */
class ArgumentImpl implements MutableArgument {

  private String name;
  private ArgumentType<?> type;
  private String description;

  public ArgumentImpl(String name, ArgumentType<?> type) {
    this.name = name;
    this.type = type;
  }

  @Override
  public MutableArgument setName(String name) {
    this.name = name;
    return this;
  }

  @Override
  public String getName() {
    return name;
  }

  @Override
  public <T> MutableArgument setType(ArgumentType<T> type) {
    this.type = type;
    return this;
  }

  @Override
  public ArgumentType<?> getType() {
    return type;
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
  public MutableArgument setDescription(String description) {
    this.description = description;
    return this;
  }

  @Override
  public String toString() {
    StringBuilder builder = new StringBuilder();
    builder.append(getClass().getName());
    builder.append(':').append(getName());
    builder.append(':').append(getType());
    if (hasDescription()) {
      builder.append(':').append(getDescription());
    }
    return builder.toString();
  }
}
