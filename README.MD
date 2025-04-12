### PULSE-OF-RAIJIN
PulseOfRaijin is a project name inspired by Japanese mythology and modern software architecture:

- **Raijin (雷神)**: The god of thunder, lightning and storms in Japanese mythology. Known for his powerful and commanding presence in the heavens.

- **Pulse**: Represents the rhythmic, consistent pattern of requests - like thunder rolling across the sky. The system handles incoming traffic with the reliability and regularity of Raijin's thunderous drums.

The name embodies our goal of building a robust, high-performance system that processes requests with the power and consistency of nature's own rhythms.

### How to use
1. Clone the repository     
2. You edit source code for your feature in cmd, internal, lib, etc.
3. You can run the application by running `make` or `make {tagname}`

### Release Process
1. Create a release branch from develop
   - Branch naming: `release/vX.Y.Z`
   - Example: `release/v1.2.0`

2. Make release preparations
   - Update version numbers
   - Update documentation
   - Fix any last minute bugs
   - No new features at this stage

3. Test thoroughly
   - Run all tests
   - Perform QA validation
   - Fix any issues found

4. Merge to main and develop
   - Merge release branch into main
   - Tag the release on main (vX.Y.Z)
   - Merge release branch back to develop
   - Delete release branch

![GitFlow Release Process](./GitFlow_Feature.png)

Note: Follow semantic versioning (MAJOR.MINOR.PATCH)
- MAJOR: Breaking changes
- MINOR: New features, backwards compatible
- PATCH: Bug fixes, backwards compatible
