package com.marshmallow.anwork.app.cli;

/**
 * This is a factory-like class that creates {@link DocumentationGenerator}'s give a
 * {@link DocumentationType}.
 *
 * <p>
 * Created Oct 16, 2017
 * </p>
 *
 * @author Andrew
 */
public class DocumentationGeneratorFactory {

  private static final DocumentationGeneratorFactory instance
      = new DocumentationGeneratorFactory();

  // Singleton pattern...
  private DocumentationGeneratorFactory() {
  }

  /**
   * Get the singleton instance {@link DocumentationGeneratorFactory}.
   *
   * @return The singleton instance {@link DocumentationGeneratorFactory}
   */
  public static DocumentationGeneratorFactory getInstance() {
    return instance;
  }

  /**
   * Create a {@link DocumentationGenerator} given a {@link DocumentationType}.
   *
   * @param type The {@link DocumentationType} of {@link DocumentationGenerator} that should be
   *     created
   * @return A {@link DocumentationGenerator} given a {@link DocumentationType}, never
   *     <code>null</code>
   */
  public DocumentationGenerator createGenerator(DocumentationType type) {
    switch (type) {
      case GITHUB_MARKDOWN:
        return new GithubReadmeDocumentationGenerator();
      default:
        throw new IllegalArgumentException("Unknown DocumentationType: " + type);
    }
  }
}
