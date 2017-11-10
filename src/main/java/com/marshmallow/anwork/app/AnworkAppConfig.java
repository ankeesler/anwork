package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.ArgumentType;
import com.marshmallow.anwork.app.cli.ArgumentValues;

import java.io.File;
import java.util.function.Consumer;

/**
 * This is a data bucket for configuration options to be passed around through
 * the ANWORK app.
 *
 * <p>
 * Created Sep 11, 2017
 * </p>
 *
 * @author Andrew
 */
public class AnworkAppConfig {

  // These are the CLI flags globally used in the ANWORK app. Each flag has a short flag
  // associated with it and maybe a {@link CliArgumentType} if it takes a parameter.
  private static enum CliFlag {
    CONTEXT("c", ArgumentType.STRING),
    PERSISTENCE_ROOT("o", ArgumentType.STRING),
    DONT_PERSIST("n"),
    DEBUG("d")
    ;

    private final String shortFlag;
    private final ArgumentType<?> argumentType;

    private CliFlag(String shortFlag, ArgumentType<?> argumentType) {
      this.shortFlag = shortFlag;
      this.argumentType = argumentType;
    }

    private CliFlag(String shortFlag) {
      // By default, flags with no parameters are translated to BOOLEAN values.
      this.shortFlag = shortFlag;
      this.argumentType = ArgumentType.BOOLEAN;
    }

    public String getShortFlag() {
      return shortFlag;
    }

    public ArgumentType<?> getArgumentType() {
      return argumentType;
    }
  }

  // These fields are set to their defaults.
  private String context = "default-context";
  private File persistenceRoot = new File(".");
  private boolean doPersist = true;
  private boolean debug = false;
  private Consumer<String> debugPrinter = new Consumer<String>() {
    @Override
    public void accept(String string) {
      if (debug) {
        System.out.println("debug: " + string);
      }
    }
  };

  /**
   * Create a configuration data object from the CLI flags to be used in the rest of the ANWORK
   * app.
   *
   * @param flags The {@link ArgumentValues} with which to initialize this object
   */
  public AnworkAppConfig(ArgumentValues flags) {
    String context = (String)getFlagValue(flags, CliFlag.CONTEXT);
    if (context != null) {
      this.context = context;
    }

    String persistenceRoot = (String)getFlagValue(flags, CliFlag.PERSISTENCE_ROOT);
    if (persistenceRoot != null) {
      this.persistenceRoot = new File(persistenceRoot);
    }

    Boolean noPersist = (Boolean)getFlagValue(flags, CliFlag.DONT_PERSIST);
    if (noPersist != null) {
      doPersist = false;
    }

    Boolean debug = (Boolean)getFlagValue(flags, CliFlag.DEBUG);
    if (debug != null && debug.equals(Boolean.TRUE)) {
      this.debug = true;
    }
  }

  public String getContext() {
    return context;
  }

  public File getPersistenceRoot() {
    return persistenceRoot;
  }

  public boolean getDoPersist() {
    return doPersist;
  }

  public Consumer<String> getDebugPrinter() {
    return debugPrinter;
  }

  private Object getFlagValue(ArgumentValues flags, CliFlag cliFlag) {
    return flags.getValue(cliFlag.getShortFlag(), cliFlag.getArgumentType());
  }
}
