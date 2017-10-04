package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.CliArgumentType;
import com.marshmallow.anwork.app.cli.CliFlags;

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

  // These fields are set to their defaults.
  private String context = "default-context";
  private File persistenceRoot;
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
   * @param flags The {@link CliFlags} with which to initialize this object
   */
  public AnworkAppConfig(CliFlags flags) {
    String context = (String)flags.getValue("c", CliArgumentType.STRING);
    if (context != null) {
      this.context = context;
    }

    String persistenceRoot = (String)flags.getValue("o", CliArgumentType.STRING);
    if (persistenceRoot != null) {
      this.persistenceRoot = new File(persistenceRoot);
    }

    Boolean noPersist = (Boolean)flags.getValue("n", CliArgumentType.BOOLEAN);
    if (noPersist != null) {
      doPersist = false;
    }

    Boolean debug = (Boolean)flags.getValue("d", CliArgumentType.BOOLEAN);
    if (debug != null && debug.equals(Boolean.TRUE)) {
      debug = true;
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

  public boolean getDebug() {
    return debug;
  }

  public Consumer<String> getDebugPrinter() {
    return debugPrinter;
  }
}
